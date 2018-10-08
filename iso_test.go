package kensho

import (
	"context"
	"testing"
)

func Test_iso3166Constraint(t *testing.T) {
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
		err := ISO3166Constraint(context.TODO(), nil, test.subject, test.arg)
		if ok := err == nil; ok != test.expected {
			t.Errorf("Expected from iso3166 constraint: %t with %T(%v)", test.expected, test.subject, test.subject)
		}
	}
}
