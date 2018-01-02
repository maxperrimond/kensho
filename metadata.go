package kensho

import (
	"fmt"
	"reflect"
	"sync"

	"strings"

	"strconv"

	"path/filepath"

	"os"

	"io/ioutil"
)

type (
	StructMetadata struct {
		StructName string
		Fields     map[string]*FieldMetadata
	}

	FieldMetadata struct {
		FieldName  string
		Validators []*ValidatorMetadata
	}

	ValidatorMetadata struct {
		Tag       string
		Validator Validator
		Arg       interface{}
	}
)

const tagName = "valid"

var metadataList = map[string]*StructMetadata{}
var metadataMutex = &sync.Mutex{}

func LoadFiles(patterns ...string) error {
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		for _, match := range matches {
			parser, ok := parsers[filepath.Ext(match)]
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
				metadataList[metadata.StructName] = metadata
			}
		}
	}

	return nil
}

func getStructMetadata(s interface{}) (*StructMetadata, error) {
	structType := reflect.TypeOf(s)
	structName := structType.String()

	metadataMutex.Lock()
	metadata, err := findOrLoadMetadata(structName, structType)
	metadataMutex.Unlock()

	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func findOrLoadMetadata(structName string, structType reflect.Type) (*StructMetadata, error) {
	if _, ok := metadataList[structName]; !ok {
		metadata, err := loadMetadataFromType(structName, structType)
		if err != nil {
			return nil, err
		}

		metadataList[structName] = metadata
	}

	return metadataList[structName], nil
}

func loadMetadataFromType(structName string, structType reflect.Type) (*StructMetadata, error) {
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
			FieldName:  fieldType.Name,
			Validators: []*ValidatorMetadata{},
		}
		metadata.Fields[fieldType.Name] = fm

		for _, validatorTag := range tags {
			c := strings.Split(strings.TrimSpace(validatorTag), "=")
			t := strings.TrimSpace(c[0])
			if t == "" {
				continue
			}

			validator, ok := validators[t]
			if !ok {
				return nil, fmt.Errorf("validator %s doesn't exists", t)
			}

			vm := &ValidatorMetadata{
				Tag:       t,
				Validator: validator,
			}

			if len(c) == 2 {
				rawArg := strings.TrimSpace(c[1])
				if rawArg == "" {
					continue
				}

				if arg, err := strconv.Atoi(rawArg); err == nil {
					vm.Arg = arg
				} else if arg, err := strconv.ParseFloat(rawArg, 64); err == nil {
					vm.Arg = arg
				} else {
					vm.Arg = strings.Trim(rawArg, `"'`)
				}
			}

			fm.Validators = append(fm.Validators, vm)
		}
	}

	return metadata, nil
}
