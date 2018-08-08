package kensho

import (
	"context"
	"testing"
)

func Test_iso3166Validator(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		subject  interface{}
		arg      interface{}
		expected bool
	}{
		{"ksdjhfjksh", nil, false},
		{"foo", "", false},
		{"", nil, true},
		{nil, nil, true},
		{"EG", "", true},
		{"EGY", "alpha3", true},
		{"818", "num", true},
	}

	for _, test := range tests {
		err := iso3166Validator(context.TODO(), nil, test.subject, test.arg)
		if ok := err == nil; ok != test.expected {
			t.Errorf("Expected from iso3166 validator: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}
