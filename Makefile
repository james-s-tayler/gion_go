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
	go run .

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## run/full: run the program against the full input file containing to produce the final Anki deck
run/full:
	go run . -input_filename="giongo-no-manga.txt"