package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
)

type Example struct {
	Giongo             string `json:"-"`
	ExampleSentence    string `json:"example_sentence" jsonschema_description:"Japanese example sentence in Kanji."`
	Hiragana           string `json:"hiragana" jsonschema_description:"The hiragana version of the example sentence separated by spaces."`
	EnglishTranslation string `json:"english_translation" jsonschema_description:"The English translation of the example sentence."`
}

type application struct {
	llm        *LLM
	giongo     []string
	examples   []Example
	outputFile *os.File
}

func New(startIndex int, inputFileName, outputFileName string) (*application, error) {
	client := openai.NewClient()
	ctx := context.Background()
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("Example"),
		Description: openai.F("Japanese example sentence, hiragana and English translation"),
		Schema:      openai.F(ExampleResponseSchema),
		Strict:      openai.Bool(true),
	}

	llm := &LLM{
		client:      client,
		ctx:         ctx,
		schemaParam: schemaParam,
	}

	output, err := os.Create(outputFileName)
	if err != nil {
		output.Close()
		return nil, fmt.Errorf("error opening file %v: %v", outputFileName, err.Error())
	}

	application := &application{
		llm:        llm,
		giongo:     []string{},
		examples:   []Example{},
		outputFile: output,
	}

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	i := 1
	for scanner.Scan() {
		if i < startIndex {
			continue
		}
		application.giongo = append(application.giongo, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return application, nil
}

func (app *application) GenerateExamples() {

	for i, v := range app.giongo {

		fmt.Printf("%d\n", i+1)
		fmt.Printf("Raw: %v\n", v)

		example, err := app.llm.GenerateExample(v)

		// need to add error handling and retries here
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("Example Sentence: %v\n", example.ExampleSentence)
		fmt.Printf("Hiragana: %v\n", example.Hiragana)
		fmt.Printf("Translation: %v\n", example.EnglishTranslation)

		app.examples = append(app.examples, example)
	}
}

func (app *application) SaveAnkiDeck() {

	for _, example := range app.examples {
		fmt.Fprintf(app.outputFile, "%v;%v;%v;%v\n",
			example.ExampleSentence,
			example.Hiragana,
			example.EnglishTranslation,
			example.Giongo,
		)
	}

	err := app.outputFile.Close()
	if err != nil {
		fmt.Printf("Error saving Anki deck: %v\n", err.Error())
		return
	}

	fmt.Println("Finished saving Anki deck.")
}
