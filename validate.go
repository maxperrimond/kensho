package kensho

import (
	"context"
	"fmt"
	"reflect"

	"strconv"

	"sync"
)

func Validate(subject interface{}) (bool, ViolationList) {
	return ValidateWithContext(context.Background(), subject)
}

func ValidateWithContext(ctx context.Context, subject interface{}) (bool, ViolationList) {
	var violations ViolationList = nil

	val := reflect.ValueOf(subject)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Type().Kind() {
	case reflect.Array, reflect.Slice:
		violations = validateList(ctx, "", val)
	case reflect.Struct:
		violations = validateStruct(ctx, "", val)
	default:
		panic(fmt.Sprintf("Cannot validate a %T, it must a struct or a list of struct", subject))
	}

	return violations == nil, violations
}

func validateStruct(ctx context.Context, path string, val reflect.Value) ViolationList {
	var wg sync.WaitGroup

	sm, err := getStructMetadata(val.Interface())
	if err != nil {
		panic(fmt.Sprintf("Unable to get validations rules for %T because %s", val.Interface(), err.Error()))
	}

	var violations ViolationList = nil
	violationsMtx := sync.Mutex{}

	for field, sfm := range sm.Fields {
		fieldVal := val.FieldByName(sfm.FieldName)
		if !fieldVal.IsValid() {
			continue
		}

		wg.Add(1)

		go func(fieldName string, fieldVal reflect.Value, metadata *FieldMetadata) {
			defer wg.Done()

			if valueViolations := validateValue(ctx, appendPath(path, fieldName), val, fieldVal, metadata); valueViolations != nil {
				violationsMtx.Lock()
				violations = append(violations, valueViolations...)
				violationsMtx.Unlock()
			}
		}(field, fieldVal, sfm)
	}

	wg.Wait()

	return violations
}

func validateList(ctx context.Context, path string, val reflect.Value) ViolationList {
	var wg sync.WaitGroup

	var violations ViolationList = nil
	violationsMtx := sync.Mutex{}

	for i := 0; i < val.Len(); i++ {
		wg.Add(1)

		go func(fieldName string, itemVal reflect.Value) {
			defer wg.Done()

			var itemViolations ViolationList = nil

			if itemVal.Kind() == reflect.Interface || itemVal.Kind() == reflect.Ptr {
				itemVal = itemVal.Elem()
			}

			switch itemVal.Type().Kind() {
			case reflect.Array, reflect.Slice:
				itemViolations = validateList(ctx, appendPath(path, fieldName), itemVal)
			case reflect.Struct:
				itemViolations = validateStruct(ctx, appendPath(path, fieldName), itemVal)
			default:
				panic(fmt.Sprintf("Cannot validate a %T, it must received a struct or a list of structs", itemVal.Interface()))
			}

			if itemViolations != nil {
				violationsMtx.Lock()
				violations = append(violations, itemViolations...)
				violationsMtx.Unlock()
			}
		}(strconv.Itoa(i), val.Index(i))
	}

	wg.Wait()

	return violations
}

func validateValue(ctx context.Context, path string, val reflect.Value, fieldVal reflect.Value, metadata *FieldMetadata) ViolationList {
	switch fieldVal.Type().Kind() {
	case reflect.Array, reflect.Slice:
		return validateArrayValue(ctx, path, val, fieldVal, metadata)
	}

	var violations ViolationList = nil

	var nested bool
	for _, vm := range metadata.Validators {
		if vm.Tag == "valid" {
			nested = true
		}

		err := vm.Validator(ctx, val.Interface(), fieldVal.Interface(), vm.Arg)
		if err != nil {
			violations = append(violations, &Violation{*err, path})
		}
	}

	if fieldVal.Kind() == reflect.Interface || fieldVal.Kind() == reflect.Ptr {
		fieldVal = fieldVal.Elem()
	}

	if nested && fieldVal.Kind() == reflect.Struct {
		structViolations := validateStruct(ctx, path, fieldVal)
		if structViolations != nil {
			violations = append(violations, structViolations...)
		}
	}

	return violations
}

func validateArrayValue(ctx context.Context, path string, val reflect.Value, fieldVal reflect.Value, metadata *FieldMetadata) ViolationList {
	var violations ViolationList = nil
	violationsMtx := sync.Mutex{}

	itemMetadata := &FieldMetadata{
		FieldName: metadata.FieldName,
	}
	for _, vm := range metadata.Validators {
		switch vm.Tag {
		case "required", "min", "max", "length":
			err := vm.Validator(ctx, val.Interface(), fieldVal.Interface(), vm.Arg)
			if err != nil {
				violations = append(violations, &Violation{*err, path})
			}
		default:
			itemMetadata.Validators = append(itemMetadata.Validators, vm)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < fieldVal.Len(); i++ {
		wg.Add(1)

		go func(fieldName string, itemVal reflect.Value) {
			defer wg.Done()

			if valueViolations := validateValue(ctx, appendPath(path, fieldName), val, itemVal, itemMetadata); valueViolations != nil {
				violationsMtx.Lock()
				violations = append(violations, valueViolations...)
				violationsMtx.Unlock()
			}
		}(strconv.Itoa(i), fieldVal.Index(i))
	}

	wg.Wait()

	return violations
}

func appendPath(path string, field string) string {
	if path == "" {
		return field
	}

	return path + "." + field
}
