package kensho

import (
	"errors"
	"fmt"
)

type (
	ValidationError struct {
		Errors []string                    `json:"errors"`
		Fields map[string]*ValidationError `json:"fields"`
	}

	ErrorBuilder func(name string, arg interface{}, previous error) string
)

var errorBuilder ErrorBuilder

var (
	missingPattern = errors.New("the pattern is missing to validate with a regex")
)

func init() {
	SetErrorBuilder(defaultErrorBuilder)
}

func SetErrorBuilder(builder ErrorBuilder) {
	errorBuilder = builder
}

func defaultErrorBuilder(name string, arg interface{}, previous error) string {
	var err string

	if previous == nil {
		err = fmt.Sprintf("not validated as %s", name)
	} else {
		err = previous.Error()
	}

	return err
}
