package kensho

import (
	"testing"
)

func TestViolation_splitPath(t *testing.T) {
	fields := (&Violation{
		Path: "foo[0].bar[102].foo",
	}).splitPath()

	if len(fields) != 5 {
		t.Error("Expected 5 fields for path \"foo[0].bar[102].foo\".")
	}
}
