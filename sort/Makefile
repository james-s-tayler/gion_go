run:
	go run .

run/threshold:
	go run . -frequency_threshold=$(threshold)

run/threshold/dataset:
	go run . -frequency_threshold=$(threshold) -frequency_data_filename=$(dataset)

run/threshold/novels:
	go run . -frequency_threshold=$(threshold) -frequency_data_filename=novels_frequency.json

run/threshold/bccw:
	go run . -frequency_threshold=$(threshold) -frequency_data_filename=bccw_frequency.json

run/threshold/netflix:
	go run . -frequency_threshold=$(threshold) -frequency_data_filename=netflix_frequency.json

.PHONY: test
test:
	go test .