package gptschema

import (
	"reflect"
	"testing"

	"github.com/akane9506/gptschema/internal"
)

func TestGenerateSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected *internal.Schema
	}{
		{
			name:     "parse common schema",
			input:    internal.CollectionWithPointers{},
			expected: &internal.CollectionWithPointersSchema,
		},
		{
			name:     "parse pointer schema",
			input:    &internal.CollectionWithPointers{},
			expected: &internal.CollectionWithPointersSchema,
		},
	}
	for _, tt := range tests {
		result, err := GenerateSchema(tt.input)
		if err != nil {
			t.Fatalf("GenerateSchema() error = %v", err)
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("expected %+v, got %+v", tt.expected, result)
		}
	}
}

func TestInvalidInputs(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		errorMsg string
	}{
		{
			name:     "nil input",
			input:    nil,
			errorMsg: "cannot generate schema for nil value",
		},
		{
			name:     "nil input",
			input:    int(3),
			errorMsg: "the schema is expected to be a Go struct",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateSchema(tt.input)
			if err == nil {
				t.Errorf("error shouldn't be nil when invalid input provided")
			}
			if err.Error() != tt.errorMsg {
				t.Errorf("mismatch error message, expect=%s, got=%s", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestGenerateSchema_WithMaxDepth(t *testing.T) {
	// Create a deeply nested structure that will exceed depth limit
	type DeepStruct struct {
		Level1 struct {
			Level2 struct {
				Level3 struct {
					Level4 struct {
						Value string `json:"value"`
					} `json:"level4"`
				} `json:"level3"`
			} `json:"level2"`
		} `json:"level1"`
	}
	tests := []struct {
		name        string
		input       interface{}
		maxDepth    int
		expectError bool
	}{
		{
			name:        "sufficient depth - should succeed",
			input:       DeepStruct{},
			maxDepth:    10,
			expectError: false,
		},
		{
			name:        "insufficient depth - should fail",
			input:       DeepStruct{},
			maxDepth:    3,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateSchema(tt.input, WithMaxDepth(tt.maxDepth))
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err != internal.ErrCircularRef {
					t.Errorf("expected ErrCircularRef, got %v", err)
				}
				if result != nil {
					t.Errorf("expected nil result on error, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected non-nil result")
				}
			}
		})
	}
}

func TestGenerateSchema_MultipleOptions(t *testing.T) {
	type Simple struct {
		Name string `json:"name"`
	}
	// Test that multiple options can be combined
	result, err := GenerateSchema(Simple{}, WithMaxDepth(10))
	if err != nil {
		t.Errorf("unexpected error with multiple options: %v", err)
	}
	if result == nil {
		t.Errorf("expected non-nil result")
	}
}
