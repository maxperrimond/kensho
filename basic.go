package kensho

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func structValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	return t.Kind() == reflect.Struct, nil
}

func stringValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	return reflect.TypeOf(value).Kind() == reflect.String, nil
}

func requiredValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return len(value.(string)) > 0, nil
	case reflect.Array, reflect.Slice, reflect.Map:
		return value != nil && reflect.ValueOf(value).Len() > 0, nil
	case reflect.Struct, reflect.Ptr, reflect.Interface:
		return value != nil, nil
	default:
		return true, nil
	}
}

func minValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	if value == nil {
		return true, nil
	}

	min, ok := arg.(int)
	if !ok {
		return false, errors.New(fmt.Sprintf("Invalid argument to min: %v", min))
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
		return false, errors.New(fmt.Sprintf("Invalid argument to max: %v", max))
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return reflect.ValueOf(value).Len() <= max, nil
	default:
		return false, nil
	}
}
