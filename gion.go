package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go" // imported as openai
)

func main() {
	fmt.Println("let's start by calling the OpenAI API and getting back an example sentence...")
	client := openai.NewClient()

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Give me an example sentence for the following: いらいら	イライラ	\"an irritating, irksome feeling (such as something stuck in the throat)\"	feelings"),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}
