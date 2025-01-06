# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/trial: run the program against a small test input file containing 3 three lines
run/trial:
	go run . -input_filename="giongo-test.txt"

## run/trial/position: run the program against a small test input file containing 3 three lines starting from $position
run/trial/position:
	go run . -input_filename="giongo-test.txt" -start_line=$(position)

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## run/full: run the program against the full input file to produce the final Anki deck
run/full:
	go run . -input_filename="giongo-no-manga.txt"

## run/full/position: run the program against the full input file starting from line $position to produce the final Anki deck
run/full/position:
	go run . -input_filename="giongo-no-manga.txt" -start_line=$(position)

