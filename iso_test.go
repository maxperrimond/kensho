package kensho

import (
	"testing"
)

func Test_iso3166Constraint(t *testing.T) {
	assertConstraintWithDataSet(t, ISO3166Constraint, []constraintCase{
		{"ksdjhfjksh", nil, false},
		{"foo", "", false},
		{"", nil, true},
		{nil, nil, true},
		{"EG", "", true},
		{"EGY", "alpha3", true},
		{"818", "num", true},
	})
}

func Test_iso639Constraint(t *testing.T) {
	assertConstraintWithDataSet(t, ISO639Constraint, []constraintCase{
		{"ksdjhfjksh", nil, false},
		{"foo", nil, false},
		{"", nil, true},
		{nil, nil, true},
		{"ja", nil, true},
		{"eng", nil, true},
		{"uz-Cyrl", nil, true},
	})
}
