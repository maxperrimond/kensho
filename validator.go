package kensho

import (
	"sync"
)

type (
	Validator struct {
		constraints map[string]Constraint
		metadata    map[string]*StructMetadata
		metadataMtx *sync.Mutex
		parsers     map[string]Parser
	}

	CustomConstraint struct {
		Name       string
		Constraint Constraint
	}
)

var defaultValidator = &Validator{
	constraints: defaultConstraints,
	metadata:    map[string]*StructMetadata{},
	metadataMtx: &sync.Mutex{},
	parsers: map[string]Parser{
		"json": jsonParse,
	},
}

func NewValidator(customConstraints ...CustomConstraint) *Validator {
	constraints := map[string]Constraint{}

	for name, constraint := range defaultConstraints {
		constraints[name] = constraint
	}
	for _, customConstraint := range customConstraints {
		constraints[customConstraint.Name] = customConstraint.Constraint
	}

	return &Validator{
		constraints: constraints,
		metadata:    map[string]*StructMetadata{},
		metadataMtx: &sync.Mutex{},
		parsers: map[string]Parser{
			"json": jsonParse,
		},
	}
}
