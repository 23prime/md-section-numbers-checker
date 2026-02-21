# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A CLI tool to check the consistency of section numbers in Markdown files.

## Commands

```bash
# Run
go run ./cmd/mdsnc examples/valid.md

# Build
go build -o bin/mdsnc ./cmd/mdsnc

# Test
go test -v ./...

# Run single test
go test -v -run TestName ./path/to/package

# Lint
golangci-lint run
```

## Architecture

- `cmd/mdsnc/` - CLI entry point
- `internal/checker/` - Core validation logic
  - `checker.go` - Validation functions and error types
  - `extractor.go` - Markdown heading parser
- `docs/spec.md` - Specification document
- `examples/` - Example Markdown files for testing

Version information is embedded at build time via `-ldflags` (see `.github/workflows/release.yml`).

## Error Codes

- `TRAILING_DOT` - Section number requires trailing dot
- `SPACING` - Exactly one space required after number
- `DEPTH_MISMATCH` - Heading level doesn't match number depth
- `MISSING_PARENT` - Parent section not defined before child
- `ORDER` - Section numbers not in ascending order
