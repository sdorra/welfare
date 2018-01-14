SRC_DIR := ./src
BIN_DIR := ./bin
SRC_FILES := $(wildcard $(SRC_DIR)/*.go)
BIN_FILES := $(patsubst $(SRC_DIR)/%.go,$(BIN_DIR)/%,$(SRC_FILES))

run: build
	docker run --rm -ti -v $(shell pwd)/$(BIN_DIR):/app -w /app ubuntu:16.04 bash

build: $(BIN_FILES)

$(BIN_DIR)/%: $(SRC_DIR)/%.go
	GOOS=linux GOARCH=amd64 go build -o $@ $<

clean:
	rm -rf $(BIN_DIR)