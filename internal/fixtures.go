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

// ==========================================

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

// ==========================================

// Struct with empty json tag (uses field name)
type StructWithEmptyTag struct {
	Name  string `json:""`
	Email string `json:"email"`
}

var StructWithEmptyTagSchema = Schema{
	"type": "object",
	"properties": Schema{
		"Name": Schema{
			"type": "string",
		},
		"email": Schema{
			"type": "string",
		},
	},
	"required":             []string{"Name", "email"},
	"additionalProperties": false,
}

// ==========================================

type NestedStruct struct {
	User   StructWithTags `json:"user"`
	Active bool           `json:"active"`
	Count  int            `json:"count,omitempty"`
}

var NestedStructSchema = Schema{
	"type": "object",
	"properties": Schema{
		"user": StructWithTagsSchema,
		"active": Schema{
			"type": "boolean",
		},
		"count": Schema{
			"type": []string{"integer", "null"},
		},
	},
	"required":             []string{"user", "active", "count"},
	"additionalProperties": false,
}

// ==========================================

// Embedded struct
type BaseInfo struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
}

type ExtendedInfo struct {
	BaseInfo        // embedded
	Title    string `json:"title"`
	Content  string `json:"content,omitempty"`
}

var ExtendedInfoSchema = Schema{
	"type": "object",
	"properties": Schema{
		"id": Schema{
			"type": "integer",
		},
		"created_at": Schema{
			"type": "string",
		},
		"title": Schema{
			"type": "string",
		},
		"content": Schema{
			"type": []string{"string", "null"},
		},
	},
	"required":             []string{"id", "created_at", "title", "content"},
	"additionalProperties": false,
}

// ==========================================

// Complex nested structures
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code,omitempty"`
}

type Company struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

type Employee struct {
	Name      string    `json:"name"`
	Companies []Company `json:"companies"`
	Tags      []string  `json:"tags,omitempty"`
}

var AddressSchema = Schema{
	"type": "object",
	"properties": Schema{
		"street": Schema{"type": "string"},
		"city":   Schema{"type": "string"},
		"zip_code": Schema{
			"type": []string{"string", "null"},
		},
	},
	"required":             []string{"street", "city", "zip_code"},
	"additionalProperties": false,
}

var CompanySchema = Schema{
	"type": "object",
	"properties": Schema{
		"name":    Schema{"type": "string"},
		"address": AddressSchema,
	},
	"required":             []string{"name", "address"},
	"additionalProperties": false,
}

var EmployeeSchema = Schema{
	"type": "object",
	"properties": Schema{
		"name": Schema{"type": "string"},
		"companies": Schema{
			"type":  "array",
			"items": CompanySchema,
		},
		"tags": Schema{
			"anyOf": []Schema{
				{
					"type": "array",
					"items": Schema{
						"type": "string",
					},
				},
				{"type": "null"},
			},
		},
	},
	"required":             []string{"name", "companies", "tags"},
	"additionalProperties": false,
}

// ==========================================

// Array of structs with pointers
type ItemWithPointer struct {
	ID   int     `json:"id"`
	Name *string `json:"name,omitempty"`
}

type CollectionWithPointers struct {
	Items []*ItemWithPointer `json:"items"`
}

var ItemWithPointerSchema = Schema{
	"type": "object",
	"properties": Schema{
		"id": Schema{"type": "integer"},
		"name": Schema{
			"type": []string{"string", "null"},
		},
	},
	"required":             []string{"id", "name"},
	"additionalProperties": false,
}

var CollectionWithPointersSchema = Schema{
	"type": "object",
	"properties": Schema{
		"items": Schema{
			"type":  "array",
			"items": ItemWithPointerSchema,
		},
	},
	"required":             []string{"items"},
	"additionalProperties": false,
}

// ==========================================

// Circular reference test (self-referencing)
type Node struct {
	Value string `json:"value"`
	Next  *Node  `json:"next,omitempty"`
}
