package kensho

import "testing"

type (
	constraintCase struct {
		value    interface{}
		arg      interface{}
		expected bool
	}
)

func assertConstraintWithDataSet(t *testing.T, constraint Constraint, dataSet []constraintCase) {
	for _, cc := range dataSet {
		assertConstraintWithContext(t, constraint, &ValidationContext{
			value: cc.value,
			arg:   cc.arg,
		}, cc.expected)
	}
}

func assertConstraintWithContext(t *testing.T, constraint Constraint, ctx *ValidationContext, expectedValid bool) {
	err := constraint(ctx)
	if err != nil {
		t.Errorf("constraint returned error: %s", err)
	}

	if expectedValid && len(ctx.ViolationList()) >= 1 {
		t.Errorf("Expected valid with value %v", ctx.value)
	} else if !expectedValid && len(ctx.ViolationList()) == 0 {
		t.Errorf("Expected invalid with value %v", ctx.value)
	}
}
