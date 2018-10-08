package kensho

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type (
	StructMetadata struct {
		StructName string
		Fields     map[string]*FieldMetadata
	}

	FieldMetadata struct {
		FieldName   string
		Constraints []*ConstraintMetadata
	}

	ConstraintMetadata struct {
		Tag        string
		Constraint Constraint
		Arg        interface{}
	}
)

const tagName = "valid"

func LoadFiles(patterns ...string) error {
	return defaultValidator.LoadFiles(patterns...)
}

func (validator *Validator) LoadFiles(patterns ...string) error {
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		for _, match := range matches {
			ext := filepath.Ext(match)
			if ext == "" {
				continue
			}

			parser, ok := validator.parsers[ext[1:]]
			if !ok {
				continue
			}

			file, err := os.Open(match)
			if err != nil {
				return err
			}

			config, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			result, err := parser(string(config))
			if err != nil {
				return err
			}

			if result == nil {
				continue
			}

			for _, metadata := range result {
				for _, field := range metadata.Fields {
					for _, val := range field.Constraints {
						constraint, ok := validator.constraints[val.Tag]
						if !ok {
							return fmt.Errorf("constraint %s doesn't exists", val.Tag)
						}

						val.Constraint = constraint
					}
				}

				validator.metadata[metadata.StructName] = metadata
			}
		}
	}

	return nil
}

func (validator *Validator) getStructMetadata(s interface{}) (*StructMetadata, error) {
	structType := reflect.TypeOf(s)
	structName := structType.String()

	validator.metadataMtx.Lock()
	metadata, err := validator.findOrLoadMetadata(structName, structType)
	validator.metadataMtx.Unlock()

	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (validator *Validator) findOrLoadMetadata(structName string, structType reflect.Type) (*StructMetadata, error) {
	if _, ok := validator.metadata[structName]; !ok {
		metadata, err := validator.loadMetadataFromType(structName, structType)
		if err != nil {
			return nil, err
		}

		validator.metadata[structName] = metadata
	}

	return validator.metadata[structName], nil
}

func (validator *Validator) loadMetadataFromType(structName string, structType reflect.Type) (*StructMetadata, error) {
	metadata := &StructMetadata{
		StructName: structName,
		Fields:     make(map[string]*FieldMetadata),
	}

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		tag := fieldType.Tag.Get(tagName)
		if tag == "" {
			continue
		}

		tags := strings.Split(tag, ",")
		if tags == nil {
			continue
		}

		fm := &FieldMetadata{
			FieldName:   fieldType.Name,
			Constraints: []*ConstraintMetadata{},
		}
		metadata.Fields[fieldType.Name] = fm

		for _, constraintTag := range tags {
			c := strings.Split(strings.TrimSpace(constraintTag), "=")
			t := strings.TrimSpace(c[0])
			if t == "" {
				continue
			}

			constraint, ok := validator.constraints[t]
			if !ok {
				return nil, fmt.Errorf("constraint %s doesn't exists", t)
			}

			cm := &ConstraintMetadata{
				Tag:        t,
				Constraint: constraint,
			}

			if len(c) == 2 {
				rawArg := strings.TrimSpace(c[1])
				if rawArg == "" {
					continue
				}

				if arg, err := strconv.Atoi(rawArg); err == nil {
					cm.Arg = arg
				} else if arg, err := strconv.ParseFloat(rawArg, 64); err == nil {
					cm.Arg = arg
				} else {
					cm.Arg = strings.Trim(rawArg, `"'`)
				}
			}

			fm.Constraints = append(fm.Constraints, cm)
		}
	}

	return metadata, nil
}
