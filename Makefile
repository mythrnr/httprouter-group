.PHONY: lint test test-json tidy
.SILENT: test-json

target ?= ./...

lint:
	golangci-lint run \
		--config=.golangci.yml \
		--print-issued-lines=false $(target)

test:
	go test -cover $(target)

test-json:
	go test -cover -json $(target)

tidy:
	go mod tidy
