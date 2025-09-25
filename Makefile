ifndef VERBOSE
MAKEFLAGS += --silent
endif

pkg ?= ./...
pwd = $(shell pwd)

.PHONY: ci-suite
ci-suite: spell-check fmt lint vulnerability-check test

.PHONY: clean
clean:
	rm -rf .cache/*

.PHONY: fmt
fmt:
	go fmt $(pkg)

.PHONY: lint
lint:
	docker pull golangci/golangci-lint:latest > /dev/null \
	&& mkdir -p .cache/golangci-lint .cache/go-build \
	&& docker run --rm \
		-v $(pwd):/app \
		-v $(pwd)/.cache:/root/.cache \
		-w /app \
		golangci/golangci-lint:latest golangci-lint run $(pkg)

.PHONY: release
release:
	if [ "$(version)" = "" ]; then \
		echo "version is required."; \
		exit 1; \
	fi \
	&& gh release create $(version) --generate-notes --target master

.PHONY: spell-check
spell-check:
	docker pull ghcr.io/streetsidesoftware/cspell:latest > /dev/null \
	&& docker run --rm \
		-v $(pwd):/workdir \
		ghcr.io/streetsidesoftware/cspell:latest \
			--config /workdir/.vscode/cspell.json "**"

.PHONY: test
test:
	go test -cover $(pkg)

.PHONY: test-json
test-json:
	go test -cover -json $(pkg)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vulnerability-check
vulnerability-check:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck -show=version ./...
