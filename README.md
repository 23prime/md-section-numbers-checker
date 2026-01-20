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

## Error Codes

| Code | Description |
| ------ | ------------- |
| `TRAILING_DOT` | Section number requires trailing dot (e.g., `1.` not `1`) |
| `SPACING` | Exactly one space required after number |
| `DEPTH_MISMATCH` | Heading level doesn't match number depth |
| `MISSING_PARENT` | Parent section not defined before child |
| `ORDER` | Section numbers not in ascending order |

## Development

### Setup

This project uses [mise](https://mise.jdx.dev/) to manage tool versions.

```bash
mise trust -q && mise install
```

### Run

- Show version

    ```bash
    go run ./cmd/md-section-numbers-checker --version
    ```

- Run for valid example

    ```bash
    go run ./cmd/md-section-numbers-checker examples/valid.md
    ```

- Run for invalid example

    ```bash
    go run ./cmd/md-section-numbers-checker examples/invalid.md
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
