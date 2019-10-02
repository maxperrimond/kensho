package kensho

import (
	"testing"
)

func Test_requiredConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, RequiredConstraint, []constraintCase{
		{nil, nil, false},
		{"", nil, false},
		{"ok", nil, true},
		{(*struct{})(nil), nil, false},
		{(interface{})(nil), nil, false},
		{&struct{}{}, nil, true},
		{struct{}{}, nil, true},
		{([]string)(nil), nil, false},
		{[]string{}, nil, true},
		{(map[string]interface{})(nil), nil, false},
		{map[string]interface{}{}, nil, true},
		{false, nil, true},
		{0, nil, true},
		{0., nil, true},
		{byte(0), nil, true},
	})
}

func Test_stringConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, StringConstraint, []constraintCase{
		{0, nil, false},
		{0., nil, false},
		{struct{}{}, nil, false},
		{true, nil, false},
		{"", nil, true},
		{nil, nil, true},
	})
}

func Test_structConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, StructConstraint, []constraintCase{
		{"", nil, false},
		{(*struct{})(nil), nil, true},
		{&struct{}{}, nil, true},
		{struct{}{}, nil, true},
		{nil, nil, true},
	})
}

func Test_lengthConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, LengthConstraint, []constraintCase{
		{"abc", -1, false},
		{"abc", 3, true},
		{"", 3, false},
		{[]string{"a", "b", "c"}, 3, true},
		{[]string{"a", "b", "c"}, 1, false},
		{[]string{}, 1, false},
	})
}

func Test_minConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, MinConstraint, []constraintCase{
		{"", 1, false},
		{"abc", 5, false},
		{"abc", 1, true},
		{[]string{"a", "b", "c"}, 3, true},
		{[]string{"a", "b", "c"}, 6, false},
		{[]string{}, 1, false},
	})
}

func Test_maxConstraint(t *testing.T) {
	assertConstraintWithDataSet(t, MaxConstraint, []constraintCase{
		{"abc", 2, false},
		{"abc", 4, true},
		{"", 1, true},
		{[]string{"a", "b", "c"}, 5, true},
		{[]string{"a", "b", "c"}, 1, false},
		{[]string{}, 1, true},
	})
}
