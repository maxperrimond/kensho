package kensho

import (
	"context"
	"testing"
)

func Test_requiredValidator(t *testing.T) {
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
		ok, _ := requiredValidator(context.Background(), nil, test.subject, nil)
		if ok != test.expected {
			t.Errorf("Expected from required validator: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}

func Test_stringValidator(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		subject  interface{}
		expected bool
	}{
		{nil, false},
		{0, false},
		{(*string)(nil), false},
		{"", true},
	}

	for _, test := range tests {
		ok, _ := stringValidator(context.Background(), nil, test.subject, nil)
		if ok != test.expected {
			t.Errorf("Expected from string validator: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}

func Test_structValidator(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		subject  interface{}
		expected bool
	}{
		{nil, false},
		{"", false},
		{(*struct{})(nil), true},
		{&struct{}{}, true},
		{struct{}{}, true},
	}

	for _, test := range tests {
		ok, _ := structValidator(context.Background(), nil, test.subject, nil)
		if ok != test.expected {
			t.Errorf("Expected from struct validator: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}
