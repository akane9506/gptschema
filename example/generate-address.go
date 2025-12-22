//go:build ignore

package example

import (
	"context"
	"encoding/json"

	"github.com/openai/openai-go/v3"
)

// A struct that will be converted to a Structured Outputs response schema
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
	schema, err := GenerateSchema(AddressItem{})

	question := "Generate a mock address for a historical russian writer"

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "address_item",
		Description: openai.String("mock address for a historical russian writ"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	print("> ")
	println(question)

	// Query the Chat Completions API
	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
		},
		Model: openai.ChatModelGPT4_1Mini2025_04_14,
	})

	if err != nil {
		panic(err.Error())
	}

	// The model responds with a JSON string, so parse it into a struct
	var address AddressItem
	// var historicalComputer HistoricalComputer
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &address)
	if err != nil {
		panic(err.Error())
	}
	println("ID:", address.ID)
	println("Name:", address.Name)
	println("Brief Intro:", address.BriefIntro)
	println("Created At:", address.CreatedAt)
	println("Updated At:", address.UpdatedAt)
	for _, tag := range address.Tags {
		println("  Tag:", tag)
	}
	println("City:", address.Address.City)
	println("Country:", address.Address.Country)
	println("Line1:", address.Address.Line1)
	println("Line2:", address.Address.Line2)
	println("Building Name:", address.Address.BuildingName)
	println("Postal Code:", address.Address.PostalCode)
	println("Region:", address.Address.Region)
}
