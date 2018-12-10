package kensho

import (
	"testing"
)

func TestViolationList_ToFormErrors(t *testing.T) {
	violationList := ViolationList{
		&Violation{
			Error{
				Error: "bad",
			},
			"",
		},
		&Violation{
			Error{
				Error: "some_error",
			},
			"foo",
		},
		&Violation{
			Error{
				Error: "some_error"},
			"foo.bar",
		},
		&Violation{
			Error{
				Error: "some_error",
			},
			"foo.0.bar.foo",
		},
	}

	formError := violationList.ToFormErrors()

	if len(formError.Errors) != 1 {
		t.Error("Expected 1 error.")
	}

	if len(formError.Fields) != 1 {
		t.Error("Expected 1 field.")
	}
}

func TestViolation_splitPath(t *testing.T) {
	fields := (&Violation{
		Path: "foo[0].bar[102].foo",
	}).splitPath()

	if len(fields) != 5 {
		t.Error("Expected 5 fields for path \"foo[0].bar[102].foo\".")
	}
}
