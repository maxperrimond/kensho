package kensho

import (
	"fmt"
	"reflect"
)

func validConstraint(_ *ValidationContext) error {
	return nil
}

func StructConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	t := reflect.TypeOf(ctx.Value())
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		ctx.BuildViolation("not_struct", nil).AddViolation()
	}

	return nil
}

func StringConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	t := reflect.TypeOf(ctx.Value())
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	if t.Kind() != reflect.String {
		ctx.BuildViolation("not_string", nil).AddViolation()
	}

	return nil
}

func RequiredConstraint(ctx *ValidationContext) error {
	if ctx.Value() != nil {
		switch reflect.TypeOf(ctx.Value()).Kind() {
		case reflect.String:
			if len(ctx.Value().(string)) > 0 {
				return nil
			}
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface:
			if !reflect.ValueOf(ctx.Value()).IsNil() {
				return nil
			}
		default:
			return nil
		}
	}

	ctx.BuildViolation("is_required", nil).AddViolation()

	return nil
}

func LengthConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	var length int
	switch ctx.Arg().(type) {
	case int:
		length = ctx.Arg().(int)
	case int64:
		length = int(ctx.Arg().(int64))
	case float64:
		length = int(ctx.Arg().(float64))
	default:
		panic(fmt.Sprintf("invalid argument to length: %v", ctx.Arg()))
	}

	switch reflect.TypeOf(ctx.Value()).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if reflect.ValueOf(ctx.Value()).Len() != length {
			ctx.BuildViolation("invalid_length", map[string]interface{}{
				"length": length,
			}).AddViolation()
		}

		return nil
	default:
		panic(fmt.Sprintf("expected a slice, map or string value"))
	}
}

func MinConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	var min int
	switch ctx.Arg().(type) {
	case int:
		min = ctx.Arg().(int)
	case int64:
		min = int(ctx.Arg().(int64))
	case float64:
		min = int(ctx.Arg().(float64))
	default:
		panic(fmt.Sprintf("invalid argument to min: %v", ctx.Arg()))
	}

	switch reflect.TypeOf(ctx.Value()).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if length := reflect.ValueOf(ctx.Value()).Len(); length < min {
			ctx.BuildViolation("too_short", map[string]interface{}{
				"min":    min,
				"length": length,
			}).AddViolation()
		}

		return nil
	default:
		panic(fmt.Sprintf("expected a slice, map or string value"))
	}
}

func MaxConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	var max int
	switch ctx.Arg().(type) {
	case int:
		max = ctx.Arg().(int)
	case int64:
		max = int(ctx.Arg().(int64))
	case float64:
		max = int(ctx.Arg().(float64))
	default:
		panic(fmt.Sprintf("invalid argument to max: %v", ctx.Arg()))
	}

	switch reflect.TypeOf(ctx.Value()).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if length := reflect.ValueOf(ctx.Value()).Len(); length > max {
			ctx.BuildViolation("too_long", map[string]interface{}{
				"max":    max,
				"length": length,
			}).AddViolation()
		}

		return nil
	default:
		panic(fmt.Sprintf("expected a slice, map or string value"))
	}
}
