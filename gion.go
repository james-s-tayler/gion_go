package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go" // imported as openai
)

type Example struct {
	ExampleSentence    string `json:"example_sentence" jsonschema_description:"Japanese example sentence."`
	Hiragana           string `json:"hiragana" jsonschema_description:"The hiragana version of the example sentence separated by spaces."`
	EnglishTranslation string `json:"english_translation" jsonschema_description:"The English translation of the example sentence."`
}

var ExampleResponseSchema = GenerateSchema[Example]()

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func main() {
	fmt.Println("let's start by calling the OpenAI API and getting back an example sentence...")
	client := openai.NewClient()
	ctx := context.Background()
	giongo := "いらいら	イライラ	\"an irritating, irksome feeling (such as something stuck in the throat)\"	feelings"
	message := "Give me an example sentence for the following: " + giongo
	fmt.Println(message)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("Example"),
		Description: openai.F("Japanese example sentence, hiragana and English translation"),
		Schema:      openai.F(ExampleResponseSchema),
		Strict:      openai.Bool(true),
	}

	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		// Only certain models can perform structured outputs
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		panic(err.Error())
	}

	example := Example{}
	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &example)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Example Sentence: %v\n", example.ExampleSentence)
	fmt.Printf("Hiranaga: %v\n", example.Hiragana)
	fmt.Printf("Translation: %v\n", example.EnglishTranslation)
}
