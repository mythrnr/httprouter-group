.PHONY: lint test test-json tidy
.SILENT: test-json

pkg ?= ./...

lint:
	golangci-lint run $(pkg)

test:
	go test -cover $(pkg)

test-json:
	go test -cover -json $(pkg)

tidy:
	go mod tidy
