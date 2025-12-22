package internal

type SimpleStruct struct {
	Name  string
	Age   int
	Email string
}

var SimpleStructSchema = Schema{
	"type": "object",
	"properties": Schema{
		"Name": Schema{
			"type": "string",
		},
		"Age": Schema{
			"type": "integer",
		},
		"Email": Schema{
			"type": "string",
		},
	},
	"required":             []string{"Name", "Age", "Email"},
	"additionalProperties": false,
}

type StructWithTags struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

var StructWithTagsSchema = Schema{
	"type": "object",
	"properties": Schema{
		"name": Schema{
			"type": "string",
		},
		"age": Schema{
			"type": "integer",
		},
		"email": Schema{
			"type": []string{"string", "null"},
		},
	},
	"required":             []string{"name", "age", "email"},
	"additionalProperties": false,
}
