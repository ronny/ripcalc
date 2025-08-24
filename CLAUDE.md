# CLAUDE.md

This file provides guidance to Claude Code (https://claude.ai/code) when working with code in this
repository.

## Project Overview

ripcalc is a CLI-only tool, written in Go, for calculating IPv4 and IPv6 address blocks.

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
