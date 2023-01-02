SHELL := /bin/bash -o pipefail

GO_PKGS ?= $(shell go list ./...)
GO_TEST_FLAGS ?= -race -v

TMP_BASE := .tmp
TMP_COVERAGE := $(TMP_BASE)/coverage

.PHONY: tools
tools:
	go install github.com/mfridman/tparse@latest

	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0

	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: lint
lint:
	go vet $(GO_PKGS)
	find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
	golangci-lint run ./...

.PHONY: test
test:
	@rm -rf $(TMP_COVERAGE)
	@mkdir -p $(TMP_COVERAGE)
	go test $(GO_TEST_FLAGS) -json -cover -coverprofile=$(TMP_COVERAGE)/coverage.txt $(GO_PKGS) | tparse
	@go tool cover -html=${TMP_COVERAGE}/coverage.txt -o $(TMP_COVERAGE)/coverage.html
	@echo
	@go tool cover -func=$(TMP_COVERAGE)/coverage.txt | grep total
	@echo
	@echo Open the coverage report
	@echo open $(TMP_COVERAGE)/coverage.html

.PYONY: dev
dev:
	go run . server -s ./tests/testdata/sample/stores -d ./tests/testdata/sample/data

.PHONY: dev_client
dev_client:
	go run . client download --down ./tests/testdata/down --watch --all