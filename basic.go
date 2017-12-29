package kensho

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

func validValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	return true, nil
}

func structValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	if value == nil {
		return false, nil
	}

	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	return t.Kind() == reflect.Struct, nil
}

func stringValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	if value == nil {
		return false, nil
	}

	return reflect.TypeOf(value).Kind() == reflect.String, nil
}

func requiredValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	if value != nil {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			if len(value.(string)) > 0 {
				return true, nil
			}
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface:
			if !reflect.ValueOf(value).IsNil() {
				return true, nil
			}
		default:
			return true, nil
		}
	}

	return false, errors.New("Is required")
}

func lengthValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	if value == nil {
		return true, nil
	}

	length, ok := arg.(int)
	if !ok {
		return false, errors.New(fmt.Sprintf("invalid argument to length: %v", length))
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return reflect.ValueOf(value).Len() == length, nil
	default:
		return false, nil
	}
}

func minValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	if value == nil {
		return true, nil
	}

	min, ok := arg.(int)
	if !ok {
		return false, errors.New(fmt.Sprintf("invalid argument to min: %v", min))
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return reflect.ValueOf(value).Len() >= min, nil
	default:
		return false, nil
	}
}

func maxValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	if value == nil {
		return true, nil
	}

	max, ok := arg.(int)
	if !ok {
		return false, errors.New(fmt.Sprintf("invalid argument to max: %v", max))
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return reflect.ValueOf(value).Len() <= max, nil
	default:
		return false, nil
	}
}
