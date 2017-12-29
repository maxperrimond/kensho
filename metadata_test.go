package kensho

import (
	"reflect"
	"testing"
)

func Test_loadMetadataFromType(t *testing.T) {
	type foo struct {
		Foo  string `valid:"required,max=30"`
		Foo2 string `valid:""`
		Foo3 string `valid:"   required,  min  = 5   ,"`
		Foo4 string `valid:"regex='^[a-z+]$'"`
	}

	fooType := reflect.TypeOf(foo{})

	metadata, err := loadMetadataFromType("foo2", fooType)
	if err != nil {
		t.Fatalf("Metadata loader failed and return this error: %s", err)
	}

	if metadata.structName != "foo2" {
		t.Errorf("Struct name is invalid, got: %s", metadata.structName)
	}

	var tests = []struct {
		fieldName  string
		validators []*validatorMetadata
	}{
		{
			"Foo",
			[]*validatorMetadata{
				{
					validatorName: "required",
					arg:           nil,
				},
				{
					validatorName: "max",
					arg:           30,
				},
			},
		},
		{
			"Foo3",
			[]*validatorMetadata{
				{
					validatorName: "required",
					arg:           nil,
				},
				{
					validatorName: "min",
					arg:           5,
				},
			},
		},
		{
			"Foo4",
			[]*validatorMetadata{
				{
					validatorName: "regex",
					arg:           "^[a-z+]$",
				},
			},
		},
	}

	for _, test := range tests {
		found := false
		for _, fm := range metadata.fields {
			if fm.fieldName == test.fieldName {
				found = true

				for _, expectedValid := range test.validators {
					vFound := false
					for _, valid := range fm.validators {
						if expectedValid.validatorName == valid.validatorName {
							vFound = true

							if expectedValid.arg != valid.arg {
								t.Errorf("Field metadata have wrong argument, expected: %T(%v), actual: %T(%v)", expectedValid.arg, expectedValid.arg, valid.arg, valid.arg)
							}

							break
						}
					}

					if !vFound {
						t.Errorf("Field validator %s not found in field metadata", expectedValid.validatorName)
					}
				}

				break
			}
		}

		if !found {
			t.Errorf("Field %s not found in metadata", test.fieldName)
		}
	}
}

func Test_loadMetadataFromType_wrongTag(t *testing.T) {
	type foo struct {
		Foo string `valid:"foo2"`
	}

	fooType := reflect.TypeOf(foo{})

	_, err := loadMetadataFromType("foo2", fooType)
	if err == nil {
		t.Error("Metadata loader should failed")
	}
}
