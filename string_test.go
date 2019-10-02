package kensho

import (
	"testing"
)

func Test_emailConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, EmailConstraint, []constraintCase{
		{"foo@goo@com", nil, false},
		{"foo@bar", nil, false},
		{"foo", nil, false},
		{"", nil, true},
		{nil, nil, true},
		{"foo@example.com", nil, true},
	})
}

func Test_uuidConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, UUIDConstraint, []constraintCase{
		{"9404926e-ef65-11e7-8c3f", nil, false},
		{"9404926eef65-11e7-8c3f-9a214cf093ae", nil, false},
		{"9404926z-ef65-11e7-8c3f-9a214cf093ae", nil, false},
		{"", nil, true},
		{nil, nil, true},
		{"9404926e-ef65-11e7-8c3f-9a214cf093ae", nil, true},
		{"dda6cd51-a791-47c6-9abf-2d835e755ad4", nil, true},
	})
}

func Test_colorHexConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, ColorHexConstraint, []constraintCase{
		{"foo", nil, false},
		{"", nil, true},
		{nil, nil, true},
		{"#fff", nil, true},
		{"#9f67cc", nil, true},
	})
}
