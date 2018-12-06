package kensho

import (
	"context"
	"fmt"
	"reflect"
)

func validConstraint(_ context.Context, _ ConstraintArgs) *Error {
	return nil
}

func StructConstraint(_ context.Context, args ConstraintArgs) *Error {
	if args.Value == nil {
		return nil
	}

	t := reflect.TypeOf(args.Value)
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

func StringConstraint(_ context.Context, args ConstraintArgs) *Error {
	if args.Value == nil {
		return nil
	}

	t := reflect.TypeOf(args.Value)
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

func RequiredConstraint(_ context.Context, args ConstraintArgs) *Error {
	if args.Value != nil {
		switch reflect.TypeOf(args.Value).Kind() {
		case reflect.String:
			if len(args.Value.(string)) > 0 {
				return nil
			}
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface:
			if !reflect.ValueOf(args.Value).IsNil() {
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

func LengthConstraint(_ context.Context, args ConstraintArgs) *Error {
	if args.Value == nil {
		return nil
	}

	var length int
	switch args.Arg.(type) {
	case int:
		length = args.Arg.(int)
	case int64:
		length = int(args.Arg.(int64))
	case float64:
		length = int(args.Arg.(float64))
	default:
		panic(fmt.Sprintf("invalid argument to length: %v", args.Arg))
	}

	switch reflect.TypeOf(args.Value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if reflect.ValueOf(args.Value).Len() != length {
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

func MinConstraint(_ context.Context, args ConstraintArgs) *Error {
	if args.Value == nil {
		return nil
	}

	var min int
	switch args.Arg.(type) {
	case int:
		min = args.Arg.(int)
	case int64:
		min = int(args.Arg.(int64))
	case float64:
		min = int(args.Arg.(float64))
	default:
		panic(fmt.Sprintf("invalid argument to min: %v", args.Arg))
	}

	switch reflect.TypeOf(args.Value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if reflect.ValueOf(args.Value).Len() < min {
			return &Error{
				Message: TranslateError("too_short", map[string]interface{}{
					"min":    min,
					"length": reflect.ValueOf(args.Value).Len(),
				}),
				Error: "too_short",
			}
		}

		return nil
	default:
		panic(fmt.Sprintf("expected a slice, map or string value"))
	}
}

func MaxConstraint(_ context.Context, args ConstraintArgs) *Error {
	if args.Value == nil {
		return nil
	}

	var max int
	switch args.Arg.(type) {
	case int:
		max = args.Arg.(int)
	case int64:
		max = int(args.Arg.(int64))
	case float64:
		max = int(args.Arg.(float64))
	default:
		panic(fmt.Sprintf("invalid argument to max: %v", args.Arg))
	}

	switch reflect.TypeOf(args.Value).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if reflect.ValueOf(args.Value).Len() > max {
			return &Error{
				Message: TranslateError("too_long", map[string]interface{}{
					"max":    max,
					"length": reflect.ValueOf(args.Value).Len(),
				}),
				Error: "too_long",
			}
		}

		return nil
	default:
		panic(fmt.Sprintf("expected a slice, map or string value"))
	}
}
