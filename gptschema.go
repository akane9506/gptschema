package gptschema

import (
	"reflect"

	"github.com/akane9506/gptschema/internal"
)

func GenerateSchema(v interface{}) (*internal.Schema, error) {
	t := reflect.TypeOf(v)
	// Dereference if pointer
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	opts := internal.DefaultOptions()
	visited := make(map[reflect.Type]bool)
	depth := 0
	result, err := internal.JsonTypeOf(t, visited, depth, opts)
	if err != nil {
		return nil, err
	}
	switch schema := result.(type) {
	case internal.Schema:
		return &schema, nil
	case map[string]interface{}:
		s := internal.Schema(schema)
		return &s, nil
	default:
		s := internal.Schema{"type": result}
		return &s, nil
	}
}
