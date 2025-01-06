package main

import (
	"math"
	"testing"
)

func TestExtractAnkiCard(t *testing.T) {
	expectedRaw := "犬がワンワンと吠えている。;いぬ が わん わん と ほえている。;The dog is barking with a woof woof sound.;わんわん	ワンワン	woof woof	onomatopoeia	dog"
	expectedHiragana := "わんわん"
	expectedKatakana := "ワンワン"
	expectedFrequency := 1
	frequency := map[string]int{
		expectedHiragana: 1,
		expectedKatakana: 2,
	}

	giongo := ExtractAnkiCard(expectedRaw, frequency)

	AssertEqual(expectedHiragana, giongo.hiragana, t)
	AssertEqual(expectedKatakana, giongo.katakana, t)
	AssertEqual(expectedRaw, giongo.raw, t)
	AssertEqual(expectedFrequency, giongo.frequency, t)
}

func TestGetFrequency(t *testing.T) {
	wanwanHiragana := "わんわん"
	wanwanKatakana := "ワンワン"
	wakuwakuHiragana := "わくわく"
	wakuwakuKatakana := "ワクワク"

	testCases := []struct {
		name              string
		expectedHiragana  string
		expectedKatakana  string
		expectedFrequency int
		frequency         map[string]int
	}{
		{
			name:             "hiragana lower than katakana",
			expectedHiragana: wanwanHiragana,
			expectedKatakana: wanwanKatakana,
			frequency: map[string]int{
				wakuwakuHiragana: 1,
				wakuwakuKatakana: 2,
				wanwanHiragana:   3,
				wanwanKatakana:   4,
			},
			expectedFrequency: 3,
		},
		{
			name:             "katakana lower than hiragana",
			expectedHiragana: wanwanHiragana,
			expectedKatakana: wanwanKatakana,
			frequency: map[string]int{
				wakuwakuHiragana: 1,
				wakuwakuKatakana: 2,
				wanwanHiragana:   4,
				wanwanKatakana:   3,
			},
			expectedFrequency: 3,
		},
		{
			name:             "hiragana not in frequency data",
			expectedHiragana: wanwanHiragana,
			expectedKatakana: wanwanKatakana,
			frequency: map[string]int{
				wakuwakuHiragana: 1,
				wakuwakuKatakana: 2,
				wanwanKatakana:   56,
			},
			expectedFrequency: 56,
		},
		{
			name:             "katakana not in frequency data",
			expectedHiragana: wanwanHiragana,
			expectedKatakana: wanwanKatakana,
			frequency: map[string]int{
				wakuwakuHiragana: 1,
				wakuwakuKatakana: 2,
				wanwanHiragana:   54,
			},
			expectedFrequency: 54,
		},
		{
			name:             "both not in frequency data",
			expectedHiragana: wanwanHiragana,
			expectedKatakana: wanwanKatakana,
			frequency: map[string]int{
				wakuwakuHiragana: 1,
				wakuwakuKatakana: 2,
			},
			expectedFrequency: math.MaxInt,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			freq := GetFrequency(testCase.expectedHiragana, testCase.expectedKatakana, testCase.frequency)
			AssertEqual(testCase.expectedFrequency, freq, t)
		})
	}
}

func TestGenerateFrequencyTag(t *testing.T) {
	testCases := []struct {
		name        string
		frequency   int
		expectedTag string
	}{
		{
			name:        "less than 1k",
			frequency:   1,
			expectedTag: "frequency-00_001-10_000",
		},
		{
			name:        "exactly 10k",
			frequency:   10_000,
			expectedTag: "frequency-00_001-10_000",
		},
		{
			name:        "exactly 10_001k",
			frequency:   10_001,
			expectedTag: "frequency-10_001-20_000",
		},
		{
			name:        "exactly 100k",
			frequency:   100_000,
			expectedTag: "frequency-90_001-100_000",
		},
		{
			name:        "100k+",
			frequency:   100_001,
			expectedTag: "frequency-100_000+",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tag := GenerateFrequencyTag(testCase.frequency)
			AssertEqual(testCase.expectedTag, tag, t)
		})
	}
}

func AssertEqual[T comparable](expected, actual T, t *testing.T) {
	t.Helper()
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
