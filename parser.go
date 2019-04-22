package kensho

import (
	"encoding/json"
)

type (
	Parser  func(config string) ([]*StructMetadata, error)
	mapping map[string]map[string][]interface{}
)

func AddParser(extension string, parser Parser) {
	defaultValidator.AddParser(extension, parser)
}

func (validator *Validator) AddParser(extension string, parser Parser) {
	validator.parsers[extension] = parser
}

func parseJSON(config string) ([]*StructMetadata, error) {
	m := mapping{}

	err := json.Unmarshal([]byte(config), &m)
	if err != nil {
		return nil, err
	}

	if len(m) == 0 {
		return nil, nil
	}

	var result []*StructMetadata

	for structName, fields := range m {
		metadata := &StructMetadata{
			StructName: structName,
			Fields:     make(map[string]*FieldMetadata),
		}

		for field, validatorList := range fields {
			fm := &FieldMetadata{
				FieldName:   field,
				Constraints: make([]*ConstraintMetadata, len(validatorList)),
			}
			metadata.Fields[field] = fm

			for i, validator := range validatorList {
				switch validator.(type) {
				case string:
					fm.Constraints[i] = &ConstraintMetadata{
						Tag: validator.(string),
					}
				case map[string]interface{}:
					config := validator.(map[string]interface{})
					for key, value := range config {
						fm.Constraints[i] = &ConstraintMetadata{
							Tag: key,
							Arg: value,
						}
						break
					}
				}
			}
		}

		if len(metadata.Fields) > 0 {
			result = append(result, metadata)
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}
