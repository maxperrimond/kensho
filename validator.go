package kensho

import "context"

type (
	Validator func(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error)
)

var validators = map[string]Validator{
	"valid":    validValidator,
	"string":   stringValidator,
	"struct":   structValidator,
	"required": requiredValidator,
	"length":   lengthValidator,
	"min":      minValidator,
	"max":      maxValidator,
	"regex":    regexValidator,
	"email":    emailValidator,
	"uuid":     uuidValidator,
	"iso3166":  iso3166Validator,
}

func AddValidator(tag string, validator Validator) {
	validators[tag] = validator
}
