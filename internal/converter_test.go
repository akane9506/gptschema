package internal

import (
	"reflect"
	"testing"

	"github.com/akane9506/jsonschema/testdata"
)

type Structured struct {
	props    map[string]any
	required []string
}

func TestJsonTypeOf(t *testing.T) {
	tests := []struct {
		name     string
		input    reflect.Type
		expected interface{}
	}{
		{
			name:     "string type",
			input:    reflect.TypeOf(""),
			expected: "string",
		},
		{
			name:  "array type",
			input: reflect.TypeOf([3]string{}),
			expected: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		{
			name:     "boolean type",
			input:    reflect.TypeOf(true),
			expected: "boolean",
		},
		{
			name:     "float type",
			input:    reflect.TypeOf(3.14),
			expected: "number",
		},
	}

	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			result := JsonTypeOf(item.input)
			if !reflect.DeepEqual(result, item.expected) {
				t.Errorf("JsonTypeOf() = %v, want %v", result, item.expected)
			}
		})
	}
}

func TestStructProperties(t *testing.T) {

	definitionStructured := Structured{
		props: map[string]any{
			"examples":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"meaning":      map[string]any{"type": "string"},
			"partOfSpeech": map[string]any{"type": "string"},
		},
		required: []string{"examples", "meaning", "partOfSpeech"},
	}

	tenseStructured := Structured{
		props: map[string]any{
			"continuous": map[string]any{"type": "string"},
			"future":     map[string]any{"type": "string"},
			"past":       map[string]any{"type": "string"},
			"present":    map[string]any{"type": "string"},
		},
		required: []string{},
	}

	wordStructed := Structured{
		props: map[string]any{
			"id":            map[string]any{"type": "string"},
			"word":          map[string]any{"type": "string"},
			"pronunciation": map[string]any{"type": "string"},
			"synonyms":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"antonyms":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"relatedTerms":  map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"definitions": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type":                 "object",
					"properties":           definitionStructured.props,
					"required":             definitionStructured.required,
					"additionalProperties": false,
				},
			},
			"tenses": map[string]any{
				"type":                 "object",
				"properties":           tenseStructured.props,
				"required":             tenseStructured.required,
				"additionalProperties": false,
			},
		},
		required: []string{"id", "word", "pronunciation", "definitions"},
	}

	tests := []struct {
		name             string
		input            reflect.Type
		expectedProps    map[string]any
		expectedRequired []string
	}{
		{
			name:             "definition type",
			input:            reflect.TypeOf(testdata.Definition{}),
			expectedProps:    definitionStructured.props,
			expectedRequired: definitionStructured.required,
		},
		{
			name:             "tense type",
			input:            reflect.TypeOf(testdata.Tenses{}),
			expectedProps:    tenseStructured.props,
			expectedRequired: tenseStructured.required,
		},
		{
			name:             "word type",
			input:            reflect.TypeOf(testdata.Word{}),
			expectedProps:    wordStructed.props,
			expectedRequired: wordStructed.required,
		},
	}

	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			props, required := StructProperties(item.input)
			if !reflect.DeepEqual(props, item.expectedProps) {
				t.Errorf(
					"StructProperties(): %v \nprops = \n%v, \nwant \n%v",
					item.name,
					props,
					item.expectedProps,
				)
			}
			if !reflect.DeepEqual(required, item.expectedRequired) {
				t.Errorf(
					"StructProperties(): %v \nrequired = \n%v, \nwant \n%v",
					item.name,
					required,
					item.expectedRequired,
				)
			}
		})
	}

	t.Run("Complete schema", func(t *testing.T) {
		output := JsonTypeOf(reflect.TypeOf(testdata.Word{}))
		required := map[string](any){
			"type":                 "object",
			"properties":           wordStructed.props,
			"required":             wordStructed.required,
			"additionalProperties": false,
		}
		if !reflect.DeepEqual(output, required) {
			t.Errorf(
				"props = \n%v, \nwant \n%v",
				output,
				required,
			)
		}
	})
}
