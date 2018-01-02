package kensho

import (
	"context"
	"testing"
)

func Test_emailValidator(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		subject  interface{}
		expected bool
	}{
		{"foo@goo@com", false},
		{"foo@bar", false},
		{"foo", false},
		{"", true},
		{nil, true},
		{"foo@example.com", true},
	}

	for _, test := range tests {
		ok, _ := emailValidator(context.TODO(), nil, test.subject, nil)
		if ok != test.expected {
			t.Errorf("Expected from email validator: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}

func Test_uuidValidator(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		subject  interface{}
		expected bool
	}{
		{"9404926e-ef65-11e7-8c3f", false},
		{"9404926eef65-11e7-8c3f-9a214cf093ae", false},
		{"9404926z-ef65-11e7-8c3f-9a214cf093ae", false},
		{"", true},
		{nil, true},
		{"9404926e-ef65-11e7-8c3f-9a214cf093ae", true},
		{"dda6cd51-a791-47c6-9abf-2d835e755ad4", true},
	}

	for _, test := range tests {
		ok, _ := uuidValidator(context.TODO(), nil, test.subject, nil)
		if ok != test.expected {
			t.Errorf("Expected from email validator: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}
