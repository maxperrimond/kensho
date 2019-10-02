package kensho

import "testing"

func TestValidationContext_ViolationList(t *testing.T) {
	ctx := &ValidationContext{
		violationList: &ViolationList{},
	}
	ctx.BuildViolation("error_1", nil).AddViolation()

	newCtx := ctx.WithValue("foo")
	newCtx.BuildViolation("error_2", nil).AddViolation()

	if len(ctx.ViolationList()) != 2 {
		t.Errorf("Expected 2 violations, found %d", len(ctx.ViolationList()))
	}
}
