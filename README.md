# gion_go

This is a set of standalone programs written in Go to produce an Anki deck to help me master onomatopoeia in Japanese.

The original data was sourced from [giongo.txt](https://github.com/Pomax/nihongoresources.com/blob/master/giongo.txt).

The frequency data was sourced from the [learnjapanese.moe Google Drive](https://drive.google.com/drive/folders/1tTdLppnqMfVC5otPlX_cs4ixlIgjv_lH)

## Process
- All manga sound effects were stripped out with `cat giongo.txt | grep -vsi "manga sound effects" > giongo-no-manga.txt`.
- `giong/gion.go` program was run to loop through `giongo-no-manga.txt` calling the OpenAI API and generate example sentences for each entry.
- `sort/main.go` program was run against the `novels_frequency.json` with a threshold of 1,000,000 (containing all the data) to sort the entire deck from most frequent to least frequent usage in novels. 
- The entire deck was then reviewed manually and any errors, duplicates, or otherwise dubious entries removed or replaced.

## The Deck

- The `anki_deck.txt` deck contains ~2k cards covering ~1k onomatopoeia.
- It is comprised of 5 fields:
  - Japanese example sentence in Kanji
  - Hiragana of the example sentence
  - English translation of the Japanese sentence
  - Onomatopoeia in Japanese + English meaning
  - Frequency range (per 10k)

## Why?

Japanese feels like it has an endless amount of onomatopoeia. I once recall seeing a dictionary of onomatopoeia that had several thousand entries. Every time a new one came up that I didn't know it eventually bothered me because of how endless it seemed. So, I decided to go on a mission to learn them all and finally get a solid handle on this fun, yet slightly vexxing, part of Japanese!