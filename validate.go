package kensho

import (
	"context"
	"fmt"
	"reflect"

	"strconv"

	"sync"

	"github.com/pkg/errors"
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
		panic(errors.Wrap(err, fmt.Sprintf("Unable to get validations rules for %T", val.Interface())))
	}

	validErr := &ValidationError{
		Fields: make(map[string]*ValidationError),
	}
	for field, sfm := range sm.fields {
		wg.Add(1)

		go func(fieldName string, fieldVal reflect.Value, metadata *fieldMetadata) {
			defer wg.Done()

			if fieldValidErr := validateValue(ctx, val, fieldVal, metadata); fieldValidErr != nil {
				validErr.Fields[fieldName] = fieldValidErr
			}
		}(field, val.FieldByIndex(sfm.sfType.Index), sfm)
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
	for i := 0; i < val.Len(); i++ {
		wg.Add(1)

		go func(fieldName string, itemVal reflect.Value) {
			defer wg.Done()

			var itemValidErr *ValidationError

			switch itemVal.Type().Kind() {
			case reflect.Array, reflect.Slice:
				itemValidErr = validateList(ctx, itemVal)
			case reflect.Struct:
				itemValidErr = validateStruct(ctx, itemVal)
			default:
				panic(fmt.Sprintf("Cannot validate a %T, it must received a struct or a list of structs", itemVal.Interface()))
			}

			if itemValidErr != nil {
				validErr.Fields[fieldName] = itemValidErr
			}
		}(strconv.Itoa(i), val.Index(i))
	}

	wg.Wait()

	if len(validErr.Fields) > 0 {
		return validErr
	}

	return nil
}

func validateValue(ctx context.Context, val reflect.Value, fieldVal reflect.Value, metadata *fieldMetadata) *ValidationError {
	validErr := &ValidationError{
		Fields: make(map[string]*ValidationError),
	}

	if fieldVal.Kind() == reflect.Interface || fieldVal.Kind() == reflect.Ptr {
		fieldVal = fieldVal.Elem()
	}

	switch fieldVal.Type().Kind() {
	case reflect.Array, reflect.Slice:
		itemMetadata := &fieldMetadata{
			fieldName: metadata.fieldName,
			sfType:    metadata.sfType,
		}

		for _, vm := range metadata.validators {
			switch vm.validatorName {
			case "required", "min", "max":
				ok, err := vm.validator(ctx, val.Interface(), fieldVal.Interface(), vm.arg)
				if !ok {
					if err == nil {
						err = errors.New(fmt.Sprintf("Not validated as %s", vm.validatorName))
					}

					validErr.Errors = append(validErr.Errors, err.Error())
				}
			default:
				itemMetadata.validators = append(itemMetadata.validators, vm)
			}
		}

		var wg sync.WaitGroup
		for i := 0; i < fieldVal.Len(); i++ {
			wg.Add(1)

			go func(fieldName string, itemVal reflect.Value) {
				defer wg.Done()

				if itemValidErr := validateValue(ctx, val, itemVal, itemMetadata); itemValidErr != nil {
					validErr.Fields[fieldName] = itemValidErr
				}
			}(strconv.Itoa(i), fieldVal.Index(i))
		}

		wg.Wait()
	case reflect.Struct:
		structValErr := validateStruct(ctx, fieldVal)
		if structValErr != nil {
			validErr.Fields = structValErr.Fields
		}
	default:
		for _, vm := range metadata.validators {
			ok, err := vm.validator(ctx, val.Interface(), fieldVal.Interface(), vm.arg)
			if !ok {
				if err == nil {
					err = errors.New(fmt.Sprintf("Not validated as %s", vm.validatorName))
				}

				validErr.Errors = append(validErr.Errors, err.Error())
			}
		}
	}

	if len(validErr.Errors) > 0 || len(validErr.Fields) > 0 {
		return validErr
	}

	return nil
}
