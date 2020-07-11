.PHONY: all test build clean

all: clean test build

build: 
	mkdir -p build
	go build -o build ./...

test:
	go test -v -coverprofile=tests/results/cover.out validate.go validate_test.go update.go update_test.go schema.go client-interface.go

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/*
	go clean ./...

container:
	podman build -t  tfld-docker-prd-local.repo.14west.io/agoracxp-cron-interface:1.14.2 .

push:
	podman push tfld-docker-prd-local.repo.14west.io/agoracxp-cron-interface:1.14.2 