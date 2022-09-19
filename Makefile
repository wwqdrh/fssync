SHELL := /bin/bash -o pipefail

GO_PKGS ?= $(shell go list ./...)
GO_TEST_FLAGS ?= -race -v

TMP_BASE := .tmp
TMP_COVERAGE := $(TMP_BASE)/coverage

.PHONY: test
test:
	go vet $(GO_PKGS)
	find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
	golangci-lint run ./...
	
	@rm -rf $(TMP_COVERAGE)
	@mkdir -p $(TMP_COVERAGE)
	go test $(GO_TEST_FLAGS) -json -cover -coverprofile=$(TMP_COVERAGE)/coverage.txt $(GO_PKGS) | tparse
	@go tool cover -html=${TMP_COVERAGE}/coverage.txt -o $(TMP_COVERAGE)/coverage.html
	@echo
	@go tool cover -func=$(TMP_COVERAGE)/coverage.txt | grep total
	@echo
	@echo Open the coverage report
	@echo open $(TMP_COVERAGE)/coverage.html