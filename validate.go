package kensho

import (
	"context"
	"fmt"
	"reflect"

	"strconv"

	"sync"
)

func Validate(subject interface{}) (bool, *ValidationError) {
	return ValidateWithContext(context.Background(), subject)
}

func ValidateWithContext(ctx context.Context, subject interface{}) (bool, *ValidationError) {
	var validErr *ValidationError

	val := reflect.ValueOf(subject)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Type().Kind() {
	case reflect.Array, reflect.Slice:
		validErr = validateList(ctx, val)
	case reflect.Struct:
		validErr = validateStruct(ctx, val)
	default:
		panic(fmt.Sprintf("Cannot validate a %T, it must a struct or a list of struct", subject))
	}

	return validErr == nil, validErr
}

func validateStruct(ctx context.Context, val reflect.Value) *ValidationError {
	var wg sync.WaitGroup

	sm, err := getStructMetadata(val.Interface())
	if err != nil {
		panic(fmt.Sprintf("Unable to get validations rules for %T because %s", val.Interface(), err.Error()))
	}

	validErr := &ValidationError{
		Fields: make(map[string]*ValidationError),
	}
	validErrMtx := sync.Mutex{}
	for field, sfm := range sm.Fields {
		fieldVal := val.FieldByName(sfm.FieldName)
		if !fieldVal.IsValid() {
			continue
		}

		wg.Add(1)

		go func(fieldName string, fieldVal reflect.Value, metadata *FieldMetadata) {
			defer wg.Done()

			if fieldValidErr := validateValue(ctx, val, fieldVal, metadata); fieldValidErr != nil {
				validErrMtx.Lock()
				validErr.Fields[fieldName] = fieldValidErr
				validErrMtx.Unlock()
			}
		}(field, fieldVal, sfm)
	}

	wg.Wait()

	if len(validErr.Fields) > 0 {
		return validErr
	}

	return nil
}

func validateList(ctx context.Context, val reflect.Value) *ValidationError {
	var wg sync.WaitGroup

	validErr := &ValidationError{
		Fields: make(map[string]*ValidationError),
	}
	validErrMtx := sync.Mutex{}
	for i := 0; i < val.Len(); i++ {
		wg.Add(1)

		go func(fieldName string, itemVal reflect.Value) {
			defer wg.Done()

			var itemValidErr *ValidationError

			if itemVal.Kind() == reflect.Interface || itemVal.Kind() == reflect.Ptr {
				itemVal = itemVal.Elem()
			}

			switch itemVal.Type().Kind() {
			case reflect.Array, reflect.Slice:
				itemValidErr = validateList(ctx, itemVal)
			case reflect.Struct:
				itemValidErr = validateStruct(ctx, itemVal)
			default:
				panic(fmt.Sprintf("Cannot validate a %T, it must received a struct or a list of structs", itemVal.Interface()))
			}

			if itemValidErr != nil {
				validErrMtx.Lock()
				validErr.Fields[fieldName] = itemValidErr
				validErrMtx.Unlock()
			}
		}(strconv.Itoa(i), val.Index(i))
	}

	wg.Wait()

	if len(validErr.Fields) > 0 {
		return validErr
	}

	return nil
}

func validateValue(ctx context.Context, val reflect.Value, fieldVal reflect.Value, metadata *FieldMetadata) *ValidationError {
	switch fieldVal.Type().Kind() {
	case reflect.Array, reflect.Slice:
		return validateArrayValue(ctx, val, fieldVal, metadata)
	}

	validErr := &ValidationError{
		Fields: make(map[string]*ValidationError),
	}

	for _, vm := range metadata.Validators {
		ok, err := vm.Validator(ctx, val.Interface(), fieldVal.Interface(), vm.Arg)
		if !ok {
			validErr.Errors = append(validErr.Errors, errorBuilder(vm.Tag, vm.Arg, err))
		}
	}

	if fieldVal.Kind() == reflect.Interface || fieldVal.Kind() == reflect.Ptr {
		fieldVal = fieldVal.Elem()
	}

	if fieldVal.Kind() == reflect.Struct {
		structValErr := validateStruct(ctx, fieldVal)
		if structValErr != nil {
			validErr.Fields = structValErr.Fields
		}
	}

	if len(validErr.Errors) > 0 || len(validErr.Fields) > 0 {
		return validErr
	}

	return nil
}

func validateArrayValue(ctx context.Context, val reflect.Value, fieldVal reflect.Value, metadata *FieldMetadata) *ValidationError {
	validErr := &ValidationError{
		Fields: make(map[string]*ValidationError),
	}
	validErrMtx := sync.Mutex{}

	itemMetadata := &FieldMetadata{
		FieldName: metadata.FieldName,
	}
	for _, vm := range metadata.Validators {
		switch vm.Tag {
		case "required", "min", "max":
			ok, err := vm.Validator(ctx, val.Interface(), fieldVal.Interface(), vm.Arg)
			if !ok {
				validErr.Errors = append(validErr.Errors, errorBuilder(vm.Tag, vm.Arg, err))
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

			if itemValidErr := validateValue(ctx, val, itemVal, itemMetadata); itemValidErr != nil {
				validErrMtx.Lock()
				validErr.Fields[fieldName] = itemValidErr
				validErrMtx.Unlock()
			}
		}(strconv.Itoa(i), fieldVal.Index(i))
	}

	wg.Wait()

	if len(validErr.Errors) > 0 || len(validErr.Fields) > 0 {
		return validErr
	}

	return nil
}
