package kensho

import (
	"context"
	"testing"
)

func Test_requiredConstraint(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		subject  interface{}
		expected bool
	}{
		{nil, false},
		{"", false},
		{"ok", true},
		{(*struct{})(nil), false},
		{(interface{})(nil), false},
		{&struct{}{}, true},
		{struct{}{}, true},
		{([]string)(nil), false},
		{[]string{}, true},
		{(map[string]interface{})(nil), false},
		{map[string]interface{}{}, true},
		{false, true},
		{0, true},
		{0., true},
		{byte(0), true},
	}

	for _, test := range tests {
		err := RequiredConstraint(context.Background(), ConstraintArgs{
			Value: test.subject,
		})
		if ok := err == nil; ok != test.expected {
			t.Errorf("Expected from required constraint: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}

func Test_stringConstraint(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		subject  interface{}
		expected bool
	}{
		{0, false},
		{0., false},
		{struct{}{}, false},
		{true, false},
		{"", true},
		{nil, true},
	}

	for _, test := range tests {
		err := StringConstraint(context.Background(), ConstraintArgs{
			Value: test.subject,
		})
		if ok := err == nil; ok != test.expected {
			t.Errorf("Expected from string constraint: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}

func Test_structConstraint(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		subject  interface{}
		expected bool
	}{
		{"", false},
		{(*struct{})(nil), true},
		{&struct{}{}, true},
		{struct{}{}, true},
		{nil, true},
	}

	for _, test := range tests {
		err := StructConstraint(context.Background(), ConstraintArgs{
			Value: test.subject,
		})
		if ok := err == nil; ok != test.expected {
			t.Errorf("Expected from struct constraint: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}
