package kensho

import (
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
		ok, violations, _ := Validate(test.subject)
		if ok != test.valid {
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %s", test.valid, ok, violations)
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
		ok, violations, _ := Validate(test.subject)
		if ok != test.valid {
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %s", test.valid, ok, violations)
		}
	}
}

func TestValidate_embeddedStruct(t *testing.T) {
	type bar struct {
		Bar string `valid:"required,min=5"`
	}

	type foo3 struct {
		Foo bar `valid:"valid"`
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
		ok, violations, _ := Validate(test.subject)
		if ok != test.valid {
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %s", test.valid, ok, violations)
		}
	}
}

func TestValidate_embeddedStructPtr(t *testing.T) {
	type bar struct {
		Bar string `valid:"required,min=5"`
	}

	type foo3 struct {
		Foo *bar `valid:"valid"`
	}

	var tests = []struct {
		subject foo3
		valid   bool
	}{
		{foo3{Foo: nil}, true},
		{foo3{Foo: &bar{Bar: "test"}}, false},
		{foo3{Foo: &bar{Bar: "okokok"}}, true},
	}

	for _, test := range tests {
		ok, violations, _ := Validate(test.subject)
		if ok != test.valid {
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %s", test.valid, ok, violations)
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
		{foo4{Foo: nil}, false},
		{foo4{Foo: []string{"test", "test"}}, false},
		{foo4{Foo: []string{"okokok"}}, true},
	}

	for _, test := range tests {
		ok, violations, _ := Validate(test.subject)
		if ok != test.valid {
			t.Errorf("Wrong validation result, expected: %t, actual: %t, with error: %s", test.valid, ok, violations)
		}
	}
}

func BenchmarkValidate(b *testing.B) {
	type someStruct struct {
		String  string        `valid:"required,min=2"`
		String2 *string       `valid:"min=10"`
		Child   []*someStruct `valid:"required,valid"`
	}

	data := make([]*someStruct, 1000)

	for i := 0; i < 1000; i++ {
		var s2 *string
		if i%2 == 0 {
			s2 = nil
		} else {
			s := "hello guys!"
			s2 = &s
		}

		data[i] = &someStruct{String: "ok", String2: s2, Child: make([]*someStruct, 1000)}

		for j := 0; j < 1000; j++ {
			data[i].Child[j] = &someStruct{String: "ok", String2: nil, Child: []*someStruct{}}
		}
	}

	b.StartTimer()
	Validate(data)
	b.StopTimer()
}
