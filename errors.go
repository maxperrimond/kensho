package kensho

import "github.com/pkg/errors"

type (
	ValidationError struct {
		Errors []string                    `json:"errors"`
		Fields map[string]*ValidationError `json:"fields"`
	}
)

var (
	missingPattern = errors.New("The pattern is missing to validate with a regex")
)
