.PHONY: all test build clean

all: clean test build


build: 
	mkdir -p build
	go build -o build -tags real ./...

test:
	go test -v -coverprofile=tests/results/cover.out -tags fake ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/*
	go clean ./...

container:
	podman build -t  quay.io/14west/trackmate-cron-interface:1.15.6 .

push:
	podman push quay.io/14west/trackmate-cron-interface:1.15.6 
