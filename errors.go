package kensho

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	missingPattern = errors.New("the pattern is missing to validate with a regex")
)

type (
	Error struct {
		Code       string                 `json:"code,omitempty"`
		Error      string                 `json:"error"`
		Message    string                 `json:"message,omitempty"`
		Parameters map[string]interface{} `json:"parameters,omitempty"`
	}

	Violation struct {
		Error
		Path string `json:"path"`
	}

	ViolationList []*Violation

	FormError struct {
		Errors []*Error              `json:"errors"`
		Fields map[string]*FormError `json:"fields,omitempty"`
	}

	ErrorTranslator func(validator ValidatorMetadata, violation Violation) string
)

func (violations *ViolationList) Error() string {
	return "" // TODO: formatted errors
}

func (violations *ViolationList) ToFormErrors() *FormError {
	if len(*violations) == 0 {
		return nil
	}

	formError := &FormError{}

	for _, violation := range *violations {
		re := regexp.MustCompile(`/^(([^\.\[]++)|\[([^\]]++)\])(.*)/`)
		re.FindAllString()
	}

	return formError
}

var errorTranslator ErrorTranslator

func init() {
	SetErrorTranslator(defaultErrorTranslator)
}

func SetErrorTranslator(builder ErrorTranslator) {
	errorTranslator = builder
}

func defaultErrorTranslator(validator ValidatorMetadata, violation Violation) string {
	if violation.Message != "" {
		return violation.Message
	}

	return fmt.Sprintf("not validated as %s", validator.Tag)
}
