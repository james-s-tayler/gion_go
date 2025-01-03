package main

import (
	"context"
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

type LLM struct {
	client      *openai.Client
	ctx         context.Context
	schemaParam openai.ResponseFormatJSONSchemaJSONSchemaParam
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

func (llm *LLM) GenerateExample(giongo string) (Example, error) {
	message := "Give me an example sentence for the following: " + giongo

	chatCompletion, err := llm.client.Chat.Completions.New(llm.ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(llm.schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		return Example{}, err
	}

	example := Example{
		Giongo: giongo,
	}
	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &example)
	if err != nil {
		return Example{}, err
	}
	return example, nil
}
