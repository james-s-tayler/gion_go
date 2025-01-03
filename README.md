# gion_go

This is a standalone program written in Go to produce an Anki deck to help me master onomatopoeia in Japanese.

The original data was sourced from [giongo.txt](https://github.com/Pomax/nihongoresources.com/blob/master/giongo.txt).

All manga sound effects were stripped out with `cat giongo.txt | grep -vsi "manga sound effects" > giongo-no-manga.txt`.

WIP
- use the OpenAI API to generate example sentences for each entry in the text file.