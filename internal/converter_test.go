package internal

import (
	"reflect"
	"testing"
)

func TestDeref(t *testing.T) {
	tests := []struct {
		name     string
		input    reflect.Type
		expected reflect.Type
	}{
		{
			name:     "non-pointer type",
			input:    reflect.TypeOf(int(0)),
			expected: reflect.TypeOf(int(0)),
		},
		{
			name:     "single pointer",
			input:    reflect.TypeOf(new(int)),
			expected: reflect.TypeOf(int(0)),
		},
		{
			name:     "double pointer",
			input:    reflect.TypeOf(new(*int)),
			expected: reflect.TypeOf(int(0)),
		},
		{
			name:     "triple pointer",
			input:    reflect.TypeOf(new(**int)),
			expected: reflect.TypeOf(int(0)),
		},
		{
			name:     "pointer to string",
			input:    reflect.TypeOf(new(string)),
			expected: reflect.TypeOf(string("")),
		},
		{
			name:     "pointer to struct",
			input:    reflect.TypeOf(new(struct{ Name string })),
			expected: reflect.TypeOf(struct{ Name string }{}),
		},
		{
			name:     "pointer to slice",
			input:    reflect.TypeOf(new([]int)),
			expected: reflect.TypeOf([]int{}),
		},
		{
			name:     "pointer to map",
			input:    reflect.TypeOf(new(map[string]int)),
			expected: reflect.TypeOf(map[string]int{}),
		},
		{
			name:     "direct slice type",
			input:    reflect.TypeOf([]string{}),
			expected: reflect.TypeOf([]string{}),
		},
		{
			name:     "direct struct type",
			input:    reflect.TypeOf(struct{ ID int }{}),
			expected: reflect.TypeOf(struct{ ID int }{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deref(tt.input)
			if result != tt.expected {
				t.Errorf("deref() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPrimitiveTypeConversion(t *testing.T) {
	tests := []struct {
		name        string
		input       reflect.Type
		output      string
		shouldError bool
	}{
		{
			name:        "string type",
			input:       reflect.TypeOf(""),
			output:      "string",
			shouldError: false,
		},
		{
			name:        "boolean type",
			input:       reflect.TypeOf(false),
			output:      "boolean",
			shouldError: false,
		},
		{
			name:        "int type",
			input:       reflect.TypeOf(int32(0)),
			output:      "integer",
			shouldError: false,
		},
		{
			name:        "uint type",
			input:       reflect.TypeOf(uint(0)),
			output:      "integer",
			shouldError: false,
		},
		{
			name:        "number type",
			input:       reflect.TypeOf(float32(0)),
			output:      "number",
			shouldError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runJsonTypeOf(tt.input)
			if tt.shouldError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.shouldError && result == nil {
				t.Errorf("expected result but got nil")
			}
		})
	}
}

func TestArrayTypeConversion(t *testing.T) {
	tests := []struct {
		name        string
		input       reflect.Type
		expected    interface{}
		shouldError bool
	}{
		{
			name:  "string items",
			input: reflect.TypeOf([]string{}),
			expected: Schema{
				"type": "string",
			},
			shouldError: false,
		},
		{
			name:  "string pointer items",
			input: reflect.TypeOf([]*string{}),
			expected: Schema{
				"type": "string",
			},
			shouldError: false,
		},
		{
			name:  "integer items",
			input: reflect.TypeOf([]int{}),
			expected: Schema{
				"type": "integer",
			},
			shouldError: false,
		},
		{
			name:  "array items",
			input: reflect.TypeOf([][]int{}),
			expected: Schema{
				"type": "array",
				"items": Schema{
					"type": "integer",
				},
			},
			shouldError: false,
		},
		{
			name:  "3D array items",
			input: reflect.TypeOf([][][]*int{}),
			expected: Schema{
				"type": "array",
				"items": Schema{
					"type": "array",
					"items": Schema{
						"type": "integer",
					},
				},
			},
			shouldError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runParseArray(tt.input)
			if tt.shouldError && err == nil {
				t.Errorf("expected error but got nil")
				return
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.shouldError {
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestStructConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    reflect.Type
		expected Schema
	}{
		{
			name:     "simple struct without tags",
			input:    reflect.TypeOf(SimpleStruct{}),
			expected: SimpleStructSchema,
		},
		{
			name:     "struct with json tags",
			input:    reflect.TypeOf(StructWithTags{}),
			expected: StructWithTagsSchema,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runJsonTypeOf(tt.input)
			if err != nil {
				t.Errorf("expected error but got nil")
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
