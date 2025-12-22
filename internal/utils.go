package internal

import (
	"reflect"
)

// Test Utils
func getInputs() (*Options, map[reflect.Type]bool, int) {
	opts := DefaultOptions()
	visited := make(map[reflect.Type]bool)
	depth := 0
	return opts, visited, depth
}

func runJsonTypeOf(input reflect.Type) (interface{}, error) {
	opts, visited, depth := getInputs()
	result, err := jsonTypeOf(input, visited, depth, opts)
	return result, err
}

func runParseArray(input reflect.Type) (interface{}, error) {
	opts, visited, depth := getInputs()
	result, err := parseArrayItemType(input, visited, depth, opts)
	return result, err
}
