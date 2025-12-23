# gptschema

[![Go Reference](https://pkg.go.dev/badge/github.com/akane9506/gptschema.svg)](https://pkg.go.dev/github.com/akane9506/gptschema)
[![Go Report Card](https://goreportcard.com/badge/github.com/akane9506/gptschema)](https://goreportcard.com/report/github.com/akane9506/gptschema)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library for generating JSON schemas from Go structs, specifically designed for [OpenAI's Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs) feature.

## Features

- 🎯 **OpenAI Optimized**: Generates schemas that comply with OpenAI's [structured output](https://platform.openai.com/docs/guides/structured-outputs) requirements
- 🔒 **Type Safe**: Leverages Go's type system to ensure schema correctness
- 🏷️ **JSON Tag Support**: Respects `json` struct tags including `omitempty` for optional fields
- 🔄 **Nested Structures**: Handles deeply nested structs, slices, and arrays
- 🛡️ **Circular Reference Detection**: Prevents infinite recursion with configurable depth limits
- 🎨 **Embedded Structs**: Automatically merges embedded struct fields

## Installation

```bash
go get github.com/akane9506/gptschema
```

## Quick start
```go
package main

import (
    "fmt"
    "log"
    
    "github.com/akane9506/gptschema"
)

type Address struct {
    City       string `json:"city"`
    Country    string `json:"country"`
    PostalCode string `json:"postalCode,omitempty"`
}

func main() {
    schema, err := gptschema.GenerateSchema(Address{})
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("%+v\n", schema)
}
```

## Usage with OpenAI
```go
package main

import (
    "context"
    "encoding/json"
    "log"
    
    "github.com/akane9506/gptschema"
    "github.com/openai/openai-go/v3"
)

type AddressItem struct {
    ID         string   `json:"id"`
    Name       string   `json:"name"`
    BriefIntro string   `json:"briefIntro"`
    CreatedAt  int64    `json:"createdAt"`
    UpdatedAt  int64    `json:"updatedAt"`
    Tags       []string `json:"tags"`
    Address    Address  `json:"address"`
}

type Address struct {
    City         string `json:"city"`
    Country      string `json:"country"`
    Line1        string `json:"line1"`
    Line2        string `json:"line2,omitempty"`
    BuildingName string `json:"buildingName,omitempty"`
    PostalCode   string `json:"postalCode,omitempty"`
    Region       string `json:"region"`
}

func main() {
    client := openai.NewClient()
    ctx := context.Background()
    
    // Generate schema from Go struct
    schema, err := gptschema.GenerateSchema(AddressItem{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Create schema parameter for OpenAI
    schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
        Name:        "address_item",
        Description: openai.String("mock address for a historical writer"),
        Schema:      schema,
        Strict:      openai.Bool(true),
    }
    
    // Query the Chat Completions API
    chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
        Messages: []openai.ChatCompletionMessageParamUnion{
            openai.UserMessage("Generate a mock address for a historical russian writer"),
        },
        ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
            OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
        },
        Model: openai.ChatModelGPT4o,
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Parse the response
    var address AddressItem
    err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &address)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the structured output
    log.Printf("Generated address for: %s in %s, %s", 
        address.Name, address.Address.City, address.Address.Country)
}
```

## Advanced Usage
### Custom maximum depth
Control the maximum depth for nested struct traversal to prevent infinite recursion:
```go
type DeepStruct struct {
    Level1 struct {
        Level2 struct {
            Level3 struct {
                Value string `json:"value"`
            } `json:"level3"`
        } `json:"level2"`
    } `json:"level1"`
}

// Set maximum depth to 10
schema, err := gptschema.GenerateSchema(DeepStruct{}, gptschema.WithMaxDepth(10))
```

### Use pointers
The library handles pointers automatically:
```go
type Person struct {
    Name string `json:"name"`
}

// Both work the same way
schema1, _ := gptschema.GenerateSchema(Person{})
schema2, _ := gptschema.GenerateSchema(&Person{})
```

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the **MIT License** - see the LICENSE file for details.