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
