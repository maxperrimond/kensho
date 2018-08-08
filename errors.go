package kensho

import (
	"strings"
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
		Errors []*Error              `json:"errors,omitempty"`
		Fields map[string]*FormError `json:"fields,omitempty"`
	}
)

func (violations ViolationList) Error() string {
	return "" // TODO: formatted errors
}

func (violations *ViolationList) ToFormErrors() *FormError {
	if len(*violations) == 0 {
		return nil
	}

	formError := &FormError{}

	for _, violation := range *violations {
		cursor := formError

		if violation.Path != "" {
			fullPath := strings.Split(violation.Path, ".")
			for _, path := range fullPath {
				if cursor.Fields == nil {
					cursor.Fields = map[string]*FormError{}
				}

				if _, has := cursor.Fields[path]; !has {
					cursor.Fields[path] = &FormError{}
				}

				cursor = cursor.Fields[path]
			}
		}

		cursor.Errors = append(cursor.Errors, &violation.Error)
	}

	return formError
}
