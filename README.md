# md-section-numbers-checker

A CLI tool to check the consistency of section numbers in Markdown files.

## Requirements

- Go 1.25+

## Installation

```bash
go install github.com/23prime/md-section-numbers-checker@latest
```

## Usage

```bash
md-section-numbers-checker <file.md>
```

## Development

### Setup

This project uses [mise](https://mise.jdx.dev/) to manage tool versions.

```bash
mise trust -q && mise install
```

### Run

```bash
go run ./cmd/md-section-numbers-checker
```

### Build

```bash
go build -o bin/md-section-numbers-checker ./cmd/md-section-numbers-checker
```

### Test

```bash
go test -v ./...
```

### Lint

```bash
golangci-lint run
```

### Git Hooks

This project uses [Lefthook](https://github.com/evilmartians/lefthook) to run automated checks on pre-commit and pre-push.

```bash
lefthook install
```

## License

MIT
