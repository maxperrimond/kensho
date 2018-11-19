package kensho

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

func Validate(subject interface{}) (bool, ViolationList) {
	return defaultValidator.Validate(subject)
}

func ValidateWithContext(ctx context.Context, subject interface{}) (bool, ViolationList) {
	return defaultValidator.ValidateWithContext(ctx, subject)
}

func (validator *Validator) Validate(subject interface{}) (bool, ViolationList) {
	return validator.ValidateWithContext(context.Background(), subject)
}

func (validator *Validator) ValidateWithContext(ctx context.Context, subject interface{}) (bool, ViolationList) {
	var violations ViolationList = nil

	val := reflect.ValueOf(subject)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Type().Kind() {
	case reflect.Array, reflect.Slice:
		violations = validator.validateList(ctx, val, "", val)
	case reflect.Struct:
		violations = validator.validateStruct(ctx, val, "", val)
	default:
		panic(fmt.Sprintf("Cannot validate a %T, it must a struct or a list of struct", subject))
	}

	return violations == nil, violations
}

func (validator *Validator) validateStruct(ctx context.Context, root reflect.Value, path string, val reflect.Value) ViolationList {
	var wg sync.WaitGroup

	sm, err := validator.getStructMetadata(val.Interface())
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

			if valueViolations := validator.validateValue(ctx, root, appendPath(path, fieldName), val, fieldVal, metadata); valueViolations != nil {
				violationsMtx.Lock()
				violations = append(violations, valueViolations...)
				violationsMtx.Unlock()
			}
		}(field, fieldVal, sfm)
	}

	wg.Wait()

	return violations
}

func (validator *Validator) validateList(ctx context.Context, root reflect.Value, path string, val reflect.Value) ViolationList {
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
				itemViolations = validator.validateList(ctx, root, appendPath(path, fieldName), itemVal)
			case reflect.Struct:
				itemViolations = validator.validateStruct(ctx, root, appendPath(path, fieldName), itemVal)
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

func (validator *Validator) validateValue(ctx context.Context, root reflect.Value, path string, val reflect.Value, fieldVal reflect.Value, metadata *FieldMetadata) ViolationList {
	switch fieldVal.Type().Kind() {
	case reflect.Array, reflect.Slice:
		return validator.validateArrayValue(ctx, root, path, val, fieldVal, metadata)
	}

	var violations ViolationList = nil

	var nested bool
	for _, constraintMetadata := range metadata.Constraints {
		if constraintMetadata.Tag == "valid" {
			nested = true
		}

		err := constraintMetadata.Constraint(ctx, ConstraintArgs{
			Root:    root.Interface(),
			Subject: val.Interface(),
			Value:   fieldVal.Interface(),
			Arg:     constraintMetadata.Arg,
		})
		if err != nil {
			violations = append(violations, &Violation{*err, path})
		}
	}

	if fieldVal.Kind() == reflect.Interface || fieldVal.Kind() == reflect.Ptr {
		fieldVal = fieldVal.Elem()
	}

	if nested && fieldVal.Kind() == reflect.Struct {
		structViolations := validator.validateStruct(ctx, root, path, fieldVal)
		if structViolations != nil {
			violations = append(violations, structViolations...)
		}
	}

	return violations
}

func (validator *Validator) validateArrayValue(ctx context.Context, root reflect.Value, path string, val reflect.Value, fieldVal reflect.Value, metadata *FieldMetadata) ViolationList {
	var violations ViolationList = nil
	violationsMtx := sync.Mutex{}

	itemMetadata := &FieldMetadata{
		FieldName: metadata.FieldName,
	}
	for _, constraintMetadata := range metadata.Constraints {
		switch constraintMetadata.Tag {
		case "required", "min", "max", "length":
			err := constraintMetadata.Constraint(ctx, ConstraintArgs{
				Root:    root.Interface(),
				Subject: val.Interface(),
				Value:   fieldVal.Interface(),
				Arg:     constraintMetadata.Arg,
			})
			if err != nil {
				violations = append(violations, &Violation{*err, path})
			}
		default:
			itemMetadata.Constraints = append(itemMetadata.Constraints, constraintMetadata)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < fieldVal.Len(); i++ {
		wg.Add(1)

		go func(fieldName string, itemVal reflect.Value) {
			defer wg.Done()

			if valueViolations := validator.validateValue(ctx, root, appendPath(path, fieldName), val, itemVal, itemMetadata); valueViolations != nil {
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
