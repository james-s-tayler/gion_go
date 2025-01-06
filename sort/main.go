package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Card struct {
	hiragana  string
	katakana  string
	raw       string
	frequency int
}

func main() {
	frequencyDataFilename := flag.String("frequency_data_filename", "novels_frequency.json", "name of frequency data file to load")
	frequencyThreshold := flag.Int("frequency_threshold", 40_000, "ignore words after this frequency - default to 40k")
	flag.Parse()

	frequency, err := LoadFrequencyData(*frequencyDataFilename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Loaded frequency data of length %d\n", len(frequency))

	unsortedLines, err := LoadAnkiCardLines("giongo_anki_deck.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Loaded unsorted anki card data of length %d\n", len(unsortedLines))

	ankiCards := ExtractAnkiCards(unsortedLines, frequency, *frequencyThreshold)

	fmt.Printf("Extracted unsorted Anki card data of length %d\n", len(ankiCards))

	sortedKeys := SortKeysByFrequency(ankiCards)

	fmt.Printf("Sorted key data of length: %d\n", len(sortedKeys))

	numSaved, err := SaveAnkiDeck(sortedKeys, ankiCards)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Saved importable Anki deck of %d cards\n", numSaved)
}

func SaveAnkiDeck(sortedKeys []string, ankiCards map[string][]Card) (int, error) {
	file, err := os.Create("sorted_anki_deck.txt")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	cardNum := 0
	for _, key := range sortedKeys {
		for _, card := range ankiCards[key] {
			cardNum++
			fmt.Fprintf(file, "%v;%v\n", card.raw, GenerateFrequencyTag(card.frequency))
		}
	}
	return cardNum, nil
}

func GenerateFrequencyTag(frequency int) string {
	switch {
	case frequency >= 1 && frequency <= 10_000:
		return "frequency-00_001-10_000"
	case frequency >= 10_001 && frequency <= 20_000:
		return "frequency-10_001-20_000"
	case frequency >= 20_001 && frequency <= 30_000:
		return "frequency-20_001-30_000"
	case frequency >= 30_001 && frequency <= 40_000:
		return "frequency-30_001-40_000"
	case frequency >= 40_001 && frequency <= 50_000:
		return "frequency-40_001-50_000"
	case frequency >= 50_001 && frequency <= 60_000:
		return "frequency-50_001-60_000"
	case frequency >= 60_001 && frequency <= 70_000:
		return "frequency-60_001-70_000"
	case frequency >= 70_001 && frequency <= 80_000:
		return "frequency-70_001-80_000"
	case frequency >= 80_001 && frequency <= 90_000:
		return "frequency-80_001-90_000"
	case frequency >= 90_001 && frequency <= 100_000:
		return "frequency-90_001-100_000"
	default:
		return "frequency-100_000+"
	}

}

func LoadFrequencyData(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var rawData [][]interface{}
	err = json.NewDecoder(file).Decode(&rawData)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %v", err)
	}

	result := make(map[string]int)

	for _, entry := range rawData {
		// validate at least the right number of entries are present
		if len(entry) >= 3 {
			headword, ok1 := entry[0].(string)
			frequency, ok2 := entry[2].(float64)

			if ok1 && ok2 {
				result[headword] = int(frequency)
			}
		}
	}

	return result, nil
}

func LoadAnkiCardLines(filename string) ([]string, error) {
	ankiDeckFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer ankiDeckFile.Close()

	lines := []string{}

	scanner := bufio.NewScanner(ankiDeckFile)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %v", err)
	}

	return lines, nil
}

func ExtractAnkiCards(lines []string, frequency map[string]int, frequencyThreshold int) map[string][]Card {
	ankiCards := make(map[string][]Card)

	skippedNotFound := 0
	skippedTooInfrequent := 0
	for _, line := range lines {
		card := ExtractAnkiCard(line, frequency)

		switch {
		case card.frequency != math.MaxInt && card.frequency <= frequencyThreshold:
			if _, ok := ankiCards[card.hiragana]; !ok {
				ankiCards[card.hiragana] = make([]Card, 0)
			}
			ankiCards[card.hiragana] = append(ankiCards[card.hiragana], card)
		case card.frequency == math.MaxInt:
			skippedNotFound++
		case card.frequency > frequencyThreshold:
			skippedTooInfrequent++
		default:
			panic("shouldn't get here")
		}
	}
	fmt.Printf("skipped: %d cards due to not being found in the frequency data\n", skippedNotFound)
	fmt.Printf("skipped: %d cards due to frequency higher than %d\n", skippedTooInfrequent, frequencyThreshold)

	return ankiCards
}

func SortKeysByFrequency(unsortedCards map[string][]Card) []string {

	sortedCardKeys := make([]string, 0, len(unsortedCards))
	for key := range unsortedCards {
		sortedCardKeys = append(sortedCardKeys, key)
	}

	sort.SliceStable(sortedCardKeys, func(i, j int) bool {
		return unsortedCards[sortedCardKeys[i]][0].frequency < unsortedCards[sortedCardKeys[j]][0].frequency
	})

	/*for _, key := range sortedCardKeys {
		fmt.Printf("%d: %v\n", frequency[key], key)
	}*/

	return sortedCardKeys
}

func ExtractAnkiCard(line string, frequency map[string]int) Card {
	lastField := strings.Split(line, ";")[3]
	lastFieldColumns := strings.Split(lastField, "\t")
	hiragana := lastFieldColumns[0]
	katakana := lastFieldColumns[1]
	lowestFrequency := GetFrequency(hiragana, katakana, frequency)

	return Card{
		hiragana:  hiragana,
		katakana:  katakana,
		raw:       line,
		frequency: lowestFrequency,
	}
}

func GetFrequency(hiragana, katakana string, frequency map[string]int) int {
	hiraganaFrequency, hiraganaOk := frequency[hiragana]
	katakanaFrequency, katakanaOk := frequency[katakana]
	if !hiraganaOk {
		hiraganaFrequency = math.MaxInt
	}
	if !katakanaOk {
		katakanaFrequency = math.MaxInt
	}
	if hiraganaFrequency < katakanaFrequency {
		return hiraganaFrequency
	}
	return katakanaFrequency
}
