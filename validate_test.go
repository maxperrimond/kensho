package kensho

import (
	"encoding/json"
	"testing"
)

func TestValidate_struct(t *testing.T) {
	type foo struct {
		Foo string `valid:"required,min=5"`
	}

	var tests = []struct {
		subject foo
		valid   bool
	}{
		{foo{Foo: ""}, false},
		{foo{Foo: "test"}, false},
		{foo{Foo: "okokok"}, true},
	}

	for _, test := range tests {
		ok, err := Validate(test.subject)
		if ok != test.valid {
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %v", test.valid, ok, err)
		}
	}
}

func TestValidate_listOfStruct(t *testing.T) {
	type foo2 struct {
		Foo string `valid:"required,min=5"`
	}

	var tests = []struct {
		subject []foo2
		valid   bool
	}{
		{[]foo2{{Foo: ""}, {Foo: ""}}, false},
		{[]foo2{{Foo: "test"}, {Foo: "okokok"}}, false},
		{[]foo2{{Foo: "okokok"}}, true},
	}

	for _, test := range tests {
		ok, err := Validate(test.subject)
		if ok != test.valid {
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %v", test.valid, ok, err)
		}
	}
}

func TestValidate_embeddedStruct(t *testing.T) {
	type bar struct {
		Bar string `valid:"required,min=5"`
	}

	type foo3 struct {
		Foo bar `valid:"struct"`
	}

	var tests = []struct {
		subject foo3
		valid   bool
	}{
		{foo3{Foo: bar{Bar: ""}}, false},
		{foo3{Foo: bar{Bar: "test"}}, false},
		{foo3{Foo: bar{Bar: "okokok"}}, true},
	}

	for _, test := range tests {
		ok, err := Validate(test.subject)
		if ok != test.valid {
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %v", test.valid, ok, err)
		}
	}
}

func TestValidate_embeddedStringList(t *testing.T) {
	type foo4 struct {
		Foo []string `valid:"required,max=1"`
	}

	var tests = []struct {
		subject foo4
		valid   bool
	}{
		{foo4{Foo: []string{}}, false},
		{foo4{Foo: []string{"test", "test"}}, false},
		{foo4{Foo: []string{"okokok"}}, true},
	}

	for _, test := range tests {
		ok, err := Validate(test.subject)
		if ok != test.valid {
			b, _ := json.Marshal(err)
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %v", test.valid, ok, string(b))
		}
	}
}
