package kensho

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

func Validate(subject interface{}) (bool, ViolationList, error) {
	return defaultValidator.Validate(subject)
}

func ValidateWithContext(ctx context.Context, subject interface{}) (bool, ViolationList, error) {
	return defaultValidator.ValidateWithContext(ctx, subject)
}

func (validator *Validator) Validate(subject interface{}) (bool, ViolationList, error) {
	return validator.ValidateWithContext(context.Background(), subject)
}

func (validator *Validator) ValidateWithContext(ctx context.Context, subject interface{}) (bool, ViolationList, error) {
	var violations ViolationList
	var err error

	val := reflect.ValueOf(subject)
	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Type().Kind() {
	case reflect.Array, reflect.Slice:
		violations, err = validator.validateList(ctx, subject, "", val)
	case reflect.Struct:
		violations, err = validator.validateStruct(ctx, subject, "", val)
	default:
		panic(fmt.Sprintf("Cannot validate a %T, it must a struct or a list of struct", subject))
	}

	if err != nil {
		return false, nil, err
	}

	return len(violations) == 0, violations, nil
}

func (validator *Validator) validateStruct(ctx context.Context, root interface{}, path string, val reflect.Value) (violations ViolationList, err error) {
	var wg sync.WaitGroup

	sm, err := validator.getStructMetadata(val.Interface())
	if err != nil {
		return nil, fmt.Errorf("unable to get validations rules for %T because %s", val.Interface(), err.Error())
	}

	for field, sfm := range sm.Fields {
		fieldVal := val.FieldByName(sfm.FieldName)
		if !fieldVal.IsValid() {
			continue
		}

		wg.Add(1)

		go func(fieldName string, fieldVal reflect.Value, metadata *FieldMetadata) {
			defer wg.Done()

			valueViolations, valueErr := validator.validateValue(ctx, root, appendPath(path, fieldName), val, fieldVal, metadata)
			if valueErr != nil {
				err = valueErr
			}

			violations = append(violations, valueViolations...)
		}(field, fieldVal, sfm)
	}

	wg.Wait()

	return violations, err
}

func (validator *Validator) validateList(ctx context.Context, root interface{}, path string, val reflect.Value) (violations ViolationList, err error) {
	var wg sync.WaitGroup

	for i := 0; i < val.Len(); i++ {
		wg.Add(1)

		go func(fieldName string, itemVal reflect.Value) {
			defer wg.Done()

			var itemViolations ViolationList = nil
			var itemErr error

			if itemVal.Kind() == reflect.Interface {
				itemVal = itemVal.Elem()
			}

			if itemVal.Kind() == reflect.Ptr {
				itemVal = itemVal.Elem()
			}

			switch itemVal.Type().Kind() {
			case reflect.Array, reflect.Slice:
				itemViolations, itemErr = validator.validateList(ctx, root, appendIndex(path, fieldName), itemVal)
			case reflect.Struct:
				itemViolations, itemErr = validator.validateStruct(ctx, root, appendIndex(path, fieldName), itemVal)
			default:
				panic(fmt.Sprintf("Cannot validate a %T, it must received a struct or a list of structs", itemVal.Interface()))
			}

			if itemErr != nil {
				err = itemErr
			}

			violations = append(violations, itemViolations...)
		}(strconv.Itoa(i), val.Index(i))
	}

	wg.Wait()

	return violations, err
}

func (validator *Validator) validateValue(ctx context.Context, root interface{}, path string, val reflect.Value, fieldVal reflect.Value, metadata *FieldMetadata) (violations ViolationList, err error) {
	switch fieldVal.Type().Kind() {
	case reflect.Array, reflect.Slice:
		return validator.validateArrayValue(ctx, root, path, val, fieldVal, metadata)
	}

	var nested bool
	for _, constraintMetadata := range metadata.Constraints {
		if constraintMetadata.Tag == "valid" {
			nested = true
		}

		validationContext := &ValidationContext{
			ctx:     ctx,
			value:   fieldVal.Interface(),
			path:    path,
			subject: val.Interface(),
			root:    root,
			arg:     constraintMetadata.Arg,
		}

		err := constraintMetadata.Constraint(validationContext)
		if err != nil {
			return nil, err
		}

		violations = append(violations, validationContext.violationList...)
	}

	if fieldVal.Kind() == reflect.Interface {
		fieldVal = fieldVal.Elem()
	}

	if fieldVal.Kind() == reflect.Ptr {
		fieldVal = fieldVal.Elem()
	}

	if nested && fieldVal.Kind() == reflect.Struct {
		structViolations, err := validator.validateStruct(ctx, root, path, fieldVal)
		if err != nil {
			return nil, err
		}

		violations = append(violations, structViolations...)
	}

	return violations, nil
}

func (validator *Validator) validateArrayValue(ctx context.Context, root interface{}, path string, val reflect.Value, fieldVal reflect.Value, metadata *FieldMetadata) (violations ViolationList, err error) {
	itemMetadata := &FieldMetadata{
		FieldName: metadata.FieldName,
	}
	for _, constraintMetadata := range metadata.Constraints {
		switch constraintMetadata.Tag {
		case "required", "min", "max", "length":
			validationContext := &ValidationContext{
				ctx:     ctx,
				value:   fieldVal.Interface(),
				path:    path,
				subject: val.Interface(),
				root:    root,
				arg:     constraintMetadata.Arg,
			}

			err := constraintMetadata.Constraint(validationContext)
			if err != nil {
				return nil, err
			}

			violations = append(violations, validationContext.violationList...)
		default:
			itemMetadata.Constraints = append(itemMetadata.Constraints, constraintMetadata)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < fieldVal.Len(); i++ {
		wg.Add(1)

		go func(fieldName string, itemVal reflect.Value) {
			defer wg.Done()

			valueViolations, valueErr := validator.validateValue(ctx, root, appendIndex(path, fieldName), val, itemVal, itemMetadata)
			if valueErr != nil {
				err = valueErr
			}

			violations = append(violations, valueViolations...)
		}(strconv.Itoa(i), fieldVal.Index(i))
	}

	wg.Wait()

	return violations, err
}

func appendPath(path string, field string) string {
	if path == "" {
		return field
	}

	return path + "." + field
}

func appendIndex(path string, index string) string {
	return path + fmt.Sprintf("[%s]", index)
}
