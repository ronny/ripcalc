# CLAUDE.md

This file provides guidance to Claude Code (https://claude.ai/code) when working with code in this
repository.

## Project Overview

ripcalc is a CLI tool and a Go package for calculating IPv4 and IPv6 address blocks.

Features:
- IPv4 and IPv6 support
- Address start, end, usable address count, network count
- Binary representation of addresses
- Structured output: JSON

## Tech Stack

- Language: Go 1.25
- Build sytem: Go toolchain, Makefile, GitHub Actions for CI, goreleaser
- Code quality: golangci-lint

## Development command

Building:

```sh
make            # Alias for checks and ripcalc binary
make ripcalc    # Build ripcalc binary
```

Testing:

```sh
make test       # Run tests
```

Code Quality:

```sh
make checks     # Run tidy, format, generate, lint, and vet
make tidy       # Run go mod tidy
make lint       # Run linters
make vet        # Run go vet
make format     # Run formatters
make test       # Run tests
```

Code generation:

```sh
make generate   # Run go generate
```

## Development Notes

- Prefer small incremental changes
  - Add things that can be done later in TODO.md to keep us focused
- Always accompany changes with tests
  - Always use `*_test` package to write tests to ensure we're testing
    only via the public interface
  - Prefer table-driven tests
- Prompt to review and create a new commit before continuing with the next change
- Only use `log/slog` for logging, never `fmt.Println` or `log`
  - Only do logging in the CLI tool, never in the ripcalc package/library
- Always handle returned errors, never silently swallow errors, either propagate it up or log it
  with `log/slog` when returning the error is not possible (e.g. in a `defer`).
- Always wrap errors with `fmt.Errorf`, use the name of the function returning the error as the
  message prefix, e.g. `fmt.Errorf("pkg.FuncName: %w", err)`
