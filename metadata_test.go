package kensho

import (
	"reflect"
	"testing"
)

func Test_loadMetadataFromType(t *testing.T) {
	t.Parallel()

	validator := NewValidator()

	type foo struct {
		Foo  string `valid:"required,max=30"`
		Foo2 string `valid:""`
		Foo3 string `valid:"   required,  min  = 5   ,"`
		Foo4 string `valid:"regex='^[a-z+]$'"`
	}

	fooType := reflect.TypeOf(foo{})

	metadata, err := validator.loadMetadataFromType("foo2", fooType)
	if err != nil {
		t.Fatalf("Metadata loader failed and return this error: %s", err)
	}

	if metadata.StructName != "foo2" {
		t.Errorf("Struct name is invalid, got: %s", metadata.StructName)
	}

	var tests = []struct {
		fieldName  string
		validators []*ConstraintMetadata
	}{
		{
			"Foo",
			[]*ConstraintMetadata{
				{
					Tag: "required",
					Arg: nil,
				},
				{
					Tag: "max",
					Arg: 30,
				},
			},
		},
		{
			"Foo3",
			[]*ConstraintMetadata{
				{
					Tag: "required",
					Arg: nil,
				},
				{
					Tag: "min",
					Arg: 5,
				},
			},
		},
		{
			"Foo4",
			[]*ConstraintMetadata{
				{
					Tag: "regex",
					Arg: "^[a-z+]$",
				},
			},
		},
	}

	for _, test := range tests {
		found := false
		for _, fm := range metadata.Fields {
			if fm.FieldName == test.fieldName {
				found = true

				for _, expectedValid := range test.validators {
					vFound := false
					for _, valid := range fm.Constraints {
						if expectedValid.Tag == valid.Tag {
							vFound = true

							if expectedValid.Arg != valid.Arg {
								t.Errorf("Field metadata have wrong argument, expected: %T(%v), actual: %T(%v)", expectedValid.Arg, expectedValid.Arg, valid.Arg, valid.Arg)
							}

							break
						}
					}

					if !vFound {
						t.Errorf("Field validator %s not found in field metadata", expectedValid.Tag)
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
	t.Parallel()

	type foo struct {
		Foo string `valid:"foo2"`
	}

	fooType := reflect.TypeOf(foo{})

	_, err := defaultValidator.loadMetadataFromType("foo2", fooType)
	if err == nil {
		t.Error("Metadata loader should failed")
	}
}

func TestLoadFiles(t *testing.T) {
	validator := NewValidator()

	validator.LoadFiles("test/*.json")

	if len(validator.metadata) != 1 {
		t.Errorf("Nb of metadata expexted: %d, actual: %d", 1, len(validator.metadata))
	}
}
