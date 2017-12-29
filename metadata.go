package kensho

import (
	"fmt"
	"reflect"
	"sync"

	"strings"

	"strconv"

	"github.com/pkg/errors"
)

type (
	structMetadata struct {
		structName string
		sType      reflect.Type
		fields     map[string]*fieldMetadata
	}

	fieldMetadata struct {
		fieldName  string
		sfType     reflect.StructField
		validators []*validatorMetadata
	}

	validatorMetadata struct {
		validatorName string
		validator     Validator
		arg           interface{}
	}
)

const tagName = "valid"

var metadataList = map[string]*structMetadata{}
var metadataMutex = &sync.Mutex{}

func getStructMetadata(s interface{}) (*structMetadata, error) {
	structType := reflect.TypeOf(s)
	structName := structType.Name()

	metadataMutex.Lock()
	metadata, err := findOrLoadMetadata(structName, structType)
	metadataMutex.Unlock()

	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func findOrLoadMetadata(structName string, structType reflect.Type) (*structMetadata, error) {
	if _, ok := metadataList[structName]; !ok {
		metadata, err := loadMetadataFromType(structName, structType)
		if err != nil {
			return nil, err
		}

		metadataList[structName] = metadata
	}

	return metadataList[structName], nil
}

func loadMetadataFromType(structName string, structType reflect.Type) (*structMetadata, error) {
	metadata := &structMetadata{
		structName: structName,
		sType:      structType,
		fields:     make(map[string]*fieldMetadata),
	}

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		tag := fieldType.Tag.Get(tagName)
		if tag == "" {
			continue
		}

		configs := strings.Split(tag, ",")
		if configs == nil {
			continue
		}

		fm := &fieldMetadata{
			fieldName:  fieldType.Name,
			sfType:     fieldType,
			validators: []*validatorMetadata{},
		}
		metadata.fields[fieldType.Name] = fm

		for _, config := range configs {
			c := strings.Split(strings.TrimSpace(config), "=")
			name := strings.TrimSpace(c[0])
			if name == "" {
				continue
			}

			validator, ok := validators[name]
			if !ok {
				return nil, errors.New(fmt.Sprintf("Validator %s doesn't exists", name))
			}

			vm := &validatorMetadata{
				validatorName: name,
				validator:     validator,
			}

			if len(c) == 2 {
				rawArg := strings.TrimSpace(c[1])
				if rawArg == "" {
					continue
				}

				if arg, err := strconv.Atoi(rawArg); err == nil {
					vm.arg = arg
				} else if arg, err := strconv.ParseFloat(rawArg, 64); err == nil {
					vm.arg = arg
				} else {
					vm.arg = strings.Trim(rawArg, `"' `)
				}
			}

			fm.validators = append(fm.validators, vm)
		}
	}

	return metadata, nil
}

func LoadFiles(filePath ...string) error {
	return nil
}
