# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A CLI tool to check the consistency of section numbers in Markdown files.

## Commands

```bash
# Run
go run ./cmd/md-section-numbers-checker

# Build
go build -o bin/md-section-numbers-checker ./cmd/md-section-numbers-checker

# Test
go test -v ./...

# Run single test
go test -v -run TestName ./path/to/package

# Lint
golangci-lint run
```

## Architecture

- `cmd/md-section-numbers-checker/` - CLI entry point
- `internal/` - Internal packages (not importable by external code)

Version information is embedded at build time via `-ldflags` (see `.github/workflows/release.yml`).
