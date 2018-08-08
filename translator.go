package kensho

type (
	ErrorTranslator func(key string, parameters map[string]interface{}) string
)

var TranslateError ErrorTranslator = func(key string, parameters map[string]interface{}) string {
	return key
}
