package kensho

import (
	"fmt"
	"regexp"
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
	if len(violations) > 0 {
		err := fmt.Sprintf("Validation returned following error \"%s\"", violations[0].Message)

		if len(violations) > 1 {
			err += fmt.Sprintf(" with %d other errors", len(violations)-1)
		}

		return err
	}

	return ""
}

func (violations *ViolationList) ToFormErrors() *FormError {
	if len(*violations) == 0 {
		return nil
	}

	formError := &FormError{}

	for _, violation := range *violations {
		cursor := formError

		if violation.Path != "" {
			fullPath := violation.splitPath()
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

func (violation *Violation) splitPath() []string {
	r := regexp.MustCompile(`\[([0-9]*)]`)
	fields := strings.Split(violation.Path, ".")

	for i := 0; ; i++ {
		if len(fields) <= i {
			break
		}

		field := fields[i]

		loc := r.FindStringIndex(field)
		if loc == nil {
			continue
		}

		slicePath := field[(loc[0] + 1):(loc[1] - 1)]

		fields[i] = string(append([]byte(field)[:loc[0]], []byte(field)[loc[1]:]...))
		fields = append(fields[:i+1], append([]string{slicePath}, fields[i+1:]...)...)
	}

	return fields
}
