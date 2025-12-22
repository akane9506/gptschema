package internal

import (
	"errors"
	"reflect"
	"strings"
)

var (
	ErrUnsupportedType = errors.New("unsupported type for JSON schema")
	ErrCircularRef     = errors.New("circular reference detected")
)

// Schema represents a JSON schema
type Schema map[string]interface{}

// Options configures schema generation
type Options struct {
	AllowAdditionalProperty bool
	MaxDepth                int
}

// DefaultOptions returns default generation options
// for OpenAI structured output, false should always set in objects:
// https://platform.openai.com/docs/guides/structured-outputs#additionalproperties-false-must-always-be-set-in-objects
// which means we don't allow Map type at all
func DefaultOptions() *Options {
	return &Options{
		AllowAdditionalProperty: false,
		MaxDepth:                50,
	}
}

// ========== Helper functions ==========

// deref dereferences pointer types recursively to get the underlying type
func deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// parse json tag
func parseJSONTag(fieldName, tag string) (name string, optional bool) {
	if tag == "" {
		return fieldName, false
	}
	parts := strings.Split(tag, ",")
	if parts[0] != "" {
		name = parts[0]
	} else {
		name = fieldName
	}
	// check for omitempty
	for _, part := range parts[1:] {
		if part == "omitempty" {
			optional = true
			break
		}
	}
	return name, optional
}

// ========== Parsing functions ==========
// convert array into json type
func parseArrayItemType(
	t reflect.Type,
	visited map[reflect.Type]bool,
	depth int,
	opts *Options) (interface{}, error) {
	schema, err := JsonTypeOf(t.Elem(), visited, depth, opts)
	if err != nil {
		return nil, err
	}
	switch s := schema.(type) {
	case string:
		return Schema{"type": s}, nil
	case Schema:
		return s, nil
	default:
		return schema, nil
	}
}

// convert struct into json
func structProperties(
	t reflect.Type,
	visited map[reflect.Type]bool,
	depth int,
	opts *Options) (Schema, []string, error) {
	props := make(Schema)
	var required []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		// handle embedded structs
		if field.Anonymous {
			embeddedProps, embeddedRequired, err := structProperties(field.Type, visited, depth, opts)
			if err != nil {
				return nil, nil, err
			}
			// merge embedded properties
			for k, v := range embeddedProps {
				props[k] = v
			}
			required = append(required, embeddedRequired...)
			continue
		}
		// parse json tag
		jsonTag := field.Tag.Get("json")
		// The json:"-" tag tells the encoding/json package
		// to ignore this field during marshaling and unmarshaling.
		if jsonTag == "-" {
			continue
		}
		fieldName, isOptional := parseJSONTag(field.Name, jsonTag)
		// generate the schema of the field
		fieldSchema, err := JsonTypeOf(field.Type, visited, depth, opts)
		if err != nil {
			return nil, nil, err
		}
		switch v := fieldSchema.(type) {
		case string:
			if isOptional {
				// Although all fields must be required,
				// it is possible to emulate an optional parameter by using a union type with null.
				props[fieldName] = Schema{"type": []string{v, "null"}}
			} else {
				props[fieldName] = Schema{"type": v}
			}
		case Schema:
			if isOptional {
				props[fieldName] = Schema{
					"anyOf": []Schema{ // OpenAI supports anyOf key
						v,
						{"type": "null"},
					},
				}
			} else {
				props[fieldName] = v
			}
		default:
			props[fieldName] = v
		}
		// All fields must be in required array for OpenAI structured outputs
		required = append(required, fieldName)
	}
	return props, required, nil
}

// JsonTypeOf converts a Go reflect.Type to a JSON Schema representation
func JsonTypeOf(
	t reflect.Type,
	visited map[reflect.Type]bool,
	depth int,
	opts *Options) (interface{}, error) {
	// check depth to prevent infinite recursion
	if depth > opts.MaxDepth {
		return nil, ErrCircularRef
	}

	t = deref(t)

	if t.Kind() == reflect.Struct {
		if visited[t] {
			return nil, ErrCircularRef
		}
		visited[t] = true
		defer delete(visited, t)
	}

	switch t.Kind() {
	case reflect.String:
		return "string", nil
	case reflect.Bool:
		return "boolean", nil
	// numbers
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer", nil
	case reflect.Float32, reflect.Float64:
		return "number", nil
	//array items
	case reflect.Slice, reflect.Array:
		items, err := parseArrayItemType(t, visited, depth+1, opts)
		if err != nil {
			return nil, err
		}
		return Schema{"type": "array", "items": items}, nil
	// object item
	case reflect.Struct:
		props, required, err := structProperties(t, visited, depth+1, opts)
		if err != nil {
			return nil, err
		}
		schema := Schema{
			"type":                 "object",
			"properties":           props,
			"additionalProperties": opts.AllowAdditionalProperty,
		}
		if len(required) > 0 {
			schema["required"] = required
		}
		return schema, nil
	default:
		return nil, ErrUnsupportedType
	}
}
