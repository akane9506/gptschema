// Package gptschema provides utilities for generating JSON schemas from Go types.
// This package is designed to work with OpenAI's structured outputs feature,
// generating schemas that comply with OpenAI's requirements.

// The main function, GenerateSchema, converts any Go type into a JSON Schema
// that can be used as a response format for OpenAI API calls.

// Please refer to the example folder for the use of these functions

package gptschema

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/akane9506/gptschema/internal"
)

// Option is a function that modifies schema generation options.
// Options can be passed to GenerateSchema to customize behavior.
type Option func(*internal.Options)

// WithMaxDepth sets the maximum depth for nested struct traversal.
// This prevents infinite recursion in deeply nested or circular structures.
// The default maximum depth is 50.
//
// Example:
//
//	schema, err := GenerateSchema(MyStruct{}, WithMaxDepth(20))
func WithMaxDepth(depth int) Option {
	return func(opts *internal.Options) {
		opts.MaxDepth = depth
	}
}

// GenerateSchema converts a Go type into a JSON Schema compatible with OpenAI's structured outputs.
//
// The function accepts any Go value and generates a JSON Schema representation that follows
// OpenAI's structured output requirements, including:
//   - additionalProperties set to false for objects
//   - All fields included in the required array
//   - Support for optional fields using union types with null (via omitempty tag)
//
// Parameters:
//   - v: Any Go value whose type will be converted to a JSON Schema. The input type can either be
//     a struct or the pointer to a struct.
//
// Returns:
//   - *internal.Schema: A pointer to the generated JSON Schema as a map[string]interface{}.
//   - error: An error if the type is unsupported or if circular references are detected.
//
// Supported Types:
//   - Primitives: string, bool, int (all variants), uint (all variants), float32, float64
//   - Complex: struct, slice, array, pointer.
//   - Embedded structs are supported and their fields are merged into the parent
//
// Unsupported Types (IMPORTANT):
//   - map: Not allowed per OpenAI's additionalProperties requirement
//   - chan, func, interface, complex types
//
// JSON Tags:
//   - Use `json:"fieldName"` to specify the JSON property name
//   - Use `json:",omitempty"` to mark fields as optional (generates union with null)
//   - Use `json:"-"` to skip fields entirely
//
// Examples:
//
//	// Simple struct
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	schema, _ := GenerateSchema(Person{})
//
//	// With optional fields
//	type User struct {
//	    Email    string  `json:"email"`
//	    Phone    string  `json:"phone,omitempty"`
//	    Nickname *string `json:"nickname,omitempty"`
//	}
//	schema, _ := GenerateSchema(User{})
//
//	// With nested structs
//	type Order struct {
//	    ID      string  `json:"id"`
//	    Address Address `json:"address"`
//	}
//	schema, _ := GenerateSchema(Order{})
//
//	// With slices
//	type Tags struct {
//	    Items []string `json:"items"`
//	}
//	schema, _ := GenerateSchema(Tags{})
//
//	// With custom options
//	type DeepStruct struct {
//	    Level1 struct {
//	        Level2 struct {
//	            Value string `json:"value"`
//	        } `json:"level2"`
//	    } `json:"level1"`
//	}
//	schema, _ := GenerateSchema(DeepStruct{}, WithMaxDepth(10))
//
// Error Conditions:
//   - Returns ErrUnsupportedType if the type cannot be converted to JSON Schema
//   - Returns ErrCircularRef if circular references are detected (depth > 50 by default)
//
// Note: The generated schema sets additionalProperties to false by default,
// which is required for OpenAI's strict mode structured outputs.
func GenerateSchema(v interface{}, opts ...Option) (*internal.Schema, error) {
	t := reflect.TypeOf(v)
	if t == nil {
		return nil, fmt.Errorf("cannot generate schema for nil value")
	}
	// Dereference if pointer
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("the schema is expected to be a Go struct")
	}
	options := internal.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	visited := make(map[reflect.Type]bool)
	depth := 0
	result, err := internal.JsonTypeOf(t, visited, depth, options)
	if err != nil {
		return nil, err
	}
	schema, ok := result.(internal.Schema)
	if !ok {
		return nil, fmt.Errorf("unexpected schema type: expected internal.Schema, got %T", result)
	}
	return &schema, nil
}

// GenerateSchemaJSON converts a Go type into a JSON Schema string compatible with OpenAI's structured outputs.
//
// This is a convenience wrapper around GenerateSchema that marshals the resulting schema
// into a JSON string, making it ready to use directly in API calls or save to files.
//
// Parameters:
//   - v: Any Go value whose type will be converted to a JSON Schema. The input type can either be
//     a struct or the pointer to a struct.
//   - opts: Optional configuration functions to customize schema generation (e.g., WithMaxDepth).
//
// Returns:
//   - string: The JSON Schema as a JSON-encoded string.
//   - error: An error if the type is unsupported, if circular references are detected,
//     or if JSON marshaling fails.
//
// This function follows the same type support and requirements as GenerateSchema.
// See GenerateSchema documentation for details on supported types, JSON tags, and constraints.
//
// Examples:
//
//	// Simple struct
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	jsonSchema, _ := GenerateSchemaJSON(Person{})
//	// jsonSchema = `{"type":"object","properties":{"name":{"type":"string"},"age":{"type":"integer"}},"required":["name","age"],"additionalProperties":false}`
//
//	// With options
//	type DeepStruct struct {
//	    Level1 struct {
//	        Value string `json:"value"`
//	    } `json:"level1"`
//	}
//	jsonSchema, _ := GenerateSchemaJSON(DeepStruct{}, WithMaxDepth(10))
//
//	// Use in OpenAI API call
//	jsonSchema, err := GenerateSchemaJSON(MyResponseFormat{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Pass jsonSchema directly to API request
//
// Error Conditions:
//   - Returns same errors as GenerateSchema for type validation and circular references
//   - Returns error if JSON marshaling fails (rare, indicates internal schema structure issue)

func GenerateSchemaJSON(v interface{}, opts ...Option) (string, error) {
	schema, err := GenerateSchema(v, opts...)
	if err != nil {
		return "", err
	}
	parsedSchema, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema to JSON: %w", err)
	}
	return string(parsedSchema), nil
}
