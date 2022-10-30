ifndef VERBOSE
MAKEFLAGS += --silent
endif

pkg ?= ./...

.PHONY: fmt
fmt:
	go fmt $(pkg)

.PHONY: lint
lint:
	golangci-lint run $(pkg)

.PHONY: nancy
nancy:
	go list -json -m all | nancy sleuth

.PHONY: spell-check
spell-check:
	# npm install -g cspell@latest
	cspell lint --config .vscode/cspell.json ".*"; \
	cspell lint --config .vscode/cspell.json "**/.*"; \
	cspell lint --config .vscode/cspell.json ".{github,vscode}/**/*"; \
	cspell lint --config .vscode/cspell.json ".{github,vscode}/**/.*"; \
	cspell lint --config .vscode/cspell.json "**"

.PHONY: test
test:
	go test -cover $(pkg)

.PHONY: test-json
test-json:
	go test -cover -json $(pkg)

.PHONY: tidy
tidy:
	go mod tidy
