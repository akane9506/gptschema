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
