pkg ?= ./...

.PHONY: lint
.SILENT: lint
lint:
	golangci-lint run $(pkg)

.PHONY: test
.SILENT: test
test:
	go test -cover $(pkg)

.PHONY: test-json
.SILENT: test-json
test-json:
	go test -cover -json $(pkg)

.PHONY: tidy
.SILENT: tidy
tidy:
	go mod tidy
