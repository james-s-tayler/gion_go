package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	now := time.Now().UnixMilli()
	inputFileStartingLine := flag.Int("start_line", 1, "line number from which to start iterating through the input file")
	inputFileName := flag.String("input_filename", "giongo-test.txt", "name of the input file containing the giongo")
	outputFileName := flag.String("output_filename", fmt.Sprintf("giongo_anki_deck-%d.txt", now), "name of the file save the anki deck")
	flag.Parse()

	application, err := New(*inputFileStartingLine, *inputFileName, *outputFileName)
	if err != nil {
		fmt.Printf("Error while instantiating the application: %v", err.Error())
		return
	}

	application.GenerateExamples()
	application.SaveAnkiDeck()
	application.SaveFailed()
}
