package kensho

import (
	"context"
	"fmt"
	"reflect"
)

func validConstraint(_ context.Context, _ interface{}, _ interface{}, _ interface{}) *Error {
	return nil
}

func StructConstraint(_ context.Context, _ interface{}, value interface{}, _ interface{}) *Error {
	if value == nil {
		return nil
	}

	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return &Error{
			Message: TranslateError("not_struct", nil),
			Error:   "not_struct",
		}
	}

	return nil
}

func StringConstraint(_ context.Context, _ interface{}, value interface{}, arg interface{}) *Error {
	if value == nil {
		return nil
	}

	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	if t.Kind() != reflect.String {
		return &Error{
			Message: TranslateError("not_string", nil),
			Error:   "not_string",
		}
	}

	return nil
}

func RequiredConstraint(_ context.Context, _ interface{}, value interface{}, arg interface{}) *Error {
	if value != nil {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			if len(value.(string)) > 0 {
				return nil
			}
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface:
			if !reflect.ValueOf(value).IsNil() {
				return nil
			}
		default:
			return nil
		}
	}

	return &Error{
		Message: TranslateError("is_required", nil),
		Error:   "is_required",
	}
}

func LengthConstraint(_ context.Context, _ interface{}, value interface{}, arg interface{}) *Error {
	if value == nil {
		return nil
	}

	length, ok := arg.(int)
	if !ok {
		panic(fmt.Sprintf("invalid argument to length: %v", arg))
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if reflect.ValueOf(value).Len() != length {
			return &Error{
				Message: TranslateError("invalid_length", map[string]interface{}{
					"length": length,
				}),
				Error: "invalid_length",
			}
		}

		return nil
	default:
		panic(fmt.Sprintf("expected a slice, map or string value"))
	}
}

func MinConstraint(_ context.Context, _ interface{}, value interface{}, arg interface{}) *Error {
	if value == nil {
		return nil
	}

	min, ok := arg.(int)
	if !ok {
		panic(fmt.Sprintf("invalid argument to min: %v", arg))
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if reflect.ValueOf(value).Len() < min {
			return &Error{
				Message: TranslateError("too_short", map[string]interface{}{
					"min":    min,
					"length": reflect.ValueOf(value).Len(),
				}),
				Error: "too_short",
			}
		}

		return nil
	default:
		panic(fmt.Sprintf("expected a slice, map or string value"))
	}
}

func MaxConstraint(_ context.Context, _ interface{}, value interface{}, arg interface{}) *Error {
	if value == nil {
		return nil
	}

	max, ok := arg.(int)
	if !ok {
		panic(fmt.Sprintf("invalid argument to max: %v", arg))
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if reflect.ValueOf(value).Len() > max {
			return &Error{
				Message: TranslateError("too_long", map[string]interface{}{
					"max":    max,
					"length": reflect.ValueOf(value).Len(),
				}),
				Error: "too_long",
			}
		}

		return nil
	default:
		panic(fmt.Sprintf("expected a slice, map or string value"))
	}
}
