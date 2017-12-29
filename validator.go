package kensho

import "context"

type (
	Validator func(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error)
)

var validators = map[string]Validator{
	"string":   stringValidator,
	"struct":   structValidator,
	"required": requiredValidator,
	"min":      minValidator,
	"max":      maxValidator,
	"regex":    regexValidator,
	"email":    emailValidator,
}

func RegisterValidator(name string, validator Validator) {
	validators[name] = validator
}
