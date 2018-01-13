.DEFAULT_GOAL := all

.PHONY: dependencies
dependencies:
	glide install

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	gometalinter --vendor ./...

.PHONY: all
all: dependencies test
