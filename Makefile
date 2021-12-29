OUTPUT_DIR ?= build_/bin
TARGETS := ./cmd/agent
TARGETS += ./cmd/server

.PHONY: build
build:
	go fmt ./...  # Format all go files in the project dir.
	CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/ $(TARGETS)

.PHONY: install
install:
	CGO_ENABLED=0 go install $(TARGETS)

.PHONY: test
test:
	go test ./... -v -count 1
