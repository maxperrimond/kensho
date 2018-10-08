package kensho

import "context"

type (
	Constraint func(ctx context.Context, subject interface{}, value interface{}, arg interface{}) *Error
)

var defaultConstraints = map[string]Constraint{
	"valid":    validConstraint,
	"string":   StringConstraint,
	"struct":   StructConstraint,
	"required": RequiredConstraint,
	"length":   LengthConstraint,
	"min":      MinConstraint,
	"max":      MaxConstraint,
	"regex":    RegexConstraint,
	"email":    EmailConstraint,
	"uuid":     UUIDConstraint,
	"colorHex": ColorHexConstraint,
	"iso3166":  ISO3166Constraint,
	"country":  ISO3166Constraint,
	"iso639":   ISO639Constraint,
	"language": ISO639Constraint,
}

func AddConstraint(name string, constraint Constraint) {
	defaultValidator.AddConstraint(name, constraint)
}

func (validator *Validator) AddConstraint(name string, constraint Constraint) {
	validator.constraints[name] = constraint
}
