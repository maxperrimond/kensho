package kensho

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	Violation struct {
		Path       string                 `json:"path"`
		Error      string                 `json:"error"`
		Message    string                 `json:"message,omitempty"`
		Parameters map[string]interface{} `json:"parameters,omitempty"`
		Code       string                 `json:"code,omitempty"`
	}

	ViolationList []*Violation

	violationBuilder struct {
		path       *string
		error      string
		message    string
		parameters map[string]interface{}
		code       *string
		basePath   string

		onAdd func(violation Violation)
	}
)

func (violations ViolationList) String() string {
	if len(violations) > 0 {
		str := "Violations:"

		for _, v := range violations {
			str += fmt.Sprintf("\n\tError at %s %q", v.Path, v.Message)
		}

		return str
	}

	return "No violations"
}

func (violations *ViolationList) append(violation *Violation) {
	*violations = append(*violations, violation)
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

func (builder *violationBuilder) SetCode(code string) *violationBuilder {
	builder.code = &code

	return builder
}

func (builder *violationBuilder) AtPath(path string) *violationBuilder {
	builder.path = &path

	return builder
}

func (builder *violationBuilder) AddViolation() {
	violation := Violation{
		Error:      builder.error,
		Message:    builder.message,
		Parameters: builder.parameters,
		Path:       builder.basePath,
	}
	if builder.code != nil {
		violation.Code = *builder.code
	}
	if builder.path != nil {
		violation.Path = appendPath(builder.basePath, *builder.path)
	}

	builder.onAdd(violation)
}
