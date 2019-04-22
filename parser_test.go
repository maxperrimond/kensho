package kensho

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestJsonParser_Parse(t *testing.T) {
	file, err := os.Open("./test/user.json")
	if err != nil {
		panic(err)
	}

	config, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	result, err := parseJSON(string(config))
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Errorf("The result length should be %d but got %d", 1, len(result))
	}

	metadata := result[0]

	if metadata.StructName != "User" {
		t.Errorf("Struct name expected: %s, actual: %s", "User", metadata.StructName)
	}

	if len(metadata.Fields) != 3 {
		t.Errorf("Nb of fields expected: %d, actual: %d", 3, len(metadata.Fields))
	}
}

func TestJsonParser_Parse_EmptyJson(t *testing.T) {
	result, err := parseJSON(`{}`)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Errorf("The result length should be %d but got %d", 0, len(result))
	}
}

func TestJsonParser_Parse_OtherFormat(t *testing.T) {
	_, err := parseJSON(`[{"foo": "foo"}]`)
	if err == nil {
		t.Error("Should get an error")
	}
}
