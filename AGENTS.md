# AGENTS.md

## Project Overview

`httprouter-group` is a Go package that provides grouped routing with middleware inheritance for [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter). It allows hierarchical route definitions where child groups inherit parent paths and middleware chains, with zero performance overhead on the underlying router.

- **Module**: `github.com/mythrnr/httprouter-group`
- **Go Version**: 1.24+
- **Dependencies**: `julienschmidt/httprouter v1.3.0`, `stretchr/testify v1.11.1` (test only)

## Project Structure

```
.
├── doc.go                  # Package-level documentation (bilingual: English/Japanese)
├── route.go                # Route and Routes types
├── route_group.go          # Core RouteGroup implementation
├── route_group_test.go     # Main test suite
├── route_group_inner_test.go # Internal helper function tests
├── route_inner_test.go     # Internal Route type tests
├── example_test.go         # Example usage with middleware
├── Makefile                # Build, test, lint commands
├── .golangci.yaml          # Linter configuration
└── .github/workflows/
    ├── check-code.yaml         # CI: lint, spell-check, unit tests
    └── scan-vulnerabilities.yaml # Vulnerability scanning
```

## Build and Test Commands

All commands are available via `Makefile`. Use `VERBOSE=1` for detailed output.

```sh
make test               # Run unit tests with coverage
make lint               # Run golangci-lint via Docker
make spell-check        # Run cspell via Docker
make fmt                # Format code with go fmt
make tidy               # Run go mod tidy
make vulnerability-check # Run govulncheck
make ci-suite           # Full CI pipeline (spell-check -> fmt -> lint -> vulnerability-check -> test)
```

## Code Conventions

### Language

- All code comments are bilingual (English + Japanese).
- Log messages must be in English.
- Error messages must be in English.

### Style

- Indentation: tabs (enforced by `.editorconfig`).
- Line endings: LF, encoding: UTF-8.
- Max line length: 100 characters (enforced by golangci-lint).
- Insert a half-width space between full-width and half-width characters in Japanese text.
- Method chaining / fluent API pattern is used throughout.

### Linting

- Almost all golangci-lint linters are enabled. Exceptions: `depguard`, `varnamelen`, `wsl`.
- Complexity thresholds: cyclop 20, gocognit/gocyclo 20, nestif 4.
- Duplication threshold: 100 lines.
- See `.golangci.yaml` for full configuration.

### Testing

- Table-driven tests with `t.Parallel()`.
- Use `testify/assert` and `testify/require` for assertions.
- Internal tests are placed in `*_inner_test.go` files (same package).
- CI runs tests against Go 1.24, 1.25, and 1.26.

## Architecture

### Core Types

- **`RouteGroup`**: Hierarchical routing node holding a path segment, HTTP handlers, child groups, and middleware. Supports fluent API via method chaining.
- **`Route`**: Flattened route with full path, HTTP method, and middleware-wrapped handler.
- **`Routes`** (`[]*Route`): Sortable collection of routes (by path, then method).
- **`Middleware`** (`func(httprouter.Handle) httprouter.Handle`): Wrapper function for cross-cutting concerns.

### Key Behaviors

- **Path inheritance**: Child paths are concatenated to parent paths. Empty/root paths normalize to `/`.
- **Middleware chaining**: Middleware is applied in reverse registration order. Parent middleware executes before child middleware.
- **Route flattening**: `RouteGroup.Routes()` recursively traverses the hierarchy and produces a flat `Routes` slice with fully resolved paths and middleware chains.

## CI/CD

- **check-code.yaml**: Triggered on PRs, pushes to master, and manual dispatch. Runs golangci-lint, spell-check (cspell), and unit tests (Go 1.24/1.25/1.26 matrix).
- **scan-vulnerabilities.yaml**: Runs daily at 00:00 UTC and on manual dispatch. Uses `govulncheck`.
- **Dependabot**: Enabled for Go modules and GitHub Actions.
