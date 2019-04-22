package kensho

import "context"

type (
	ConstraintArgs struct {
		Root    interface{}
		Subject interface{}
		Value   interface{}
		Arg     interface{}
	}

	Constraint func(ctx context.Context, args ConstraintArgs) *Error

	CustomConstraint struct {
		Name       string
		Constraint Constraint
	}
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

func NewCustomConstraint(name string, constraint Constraint) CustomConstraint {
	return CustomConstraint{
		Name:       name,
		Constraint: constraint,
	}
}

func AddConstraint(name string, constraint Constraint) {
	defaultValidator.AddConstraint(name, constraint)
}

func (validator *Validator) AddConstraint(name string, constraint Constraint) {
	validator.constraints[name] = constraint
}
