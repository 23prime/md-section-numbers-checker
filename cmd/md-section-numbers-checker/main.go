package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/23prime/md-section-numbers-checker/internal/checker"
)

const (
	exitSuccess = 0
	exitError   = 1
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func main() {
	os.Exit(run())
}

func run() int {
	showHelp := flag.Bool("help", false, "Show help message")
	flag.BoolVar(showHelp, "h", false, "Show help message (shorthand)")

	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showHelp {
		showHelpMessage()
		return exitSuccess
	}

	if *showVersion {
		printVersion()
		return exitSuccess
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: md-section-numbers-checker <file.md> ...")
		return exitError
	}

	hasError := processPatterns(args)

	if hasError {
		return exitError
	}

	fmt.Println("All section numbers are valid.")
	return exitSuccess
}

func showHelpMessage() {
	fmt.Printf(`md-section-numbers-checker %s

USAGE:
  md-section-numbers-checker [OPTIONS] <file.md> ...

OPTIONS:
  -h, --help     Show this help message
      --version  Show version information

ERROR CODES:
  TRAILING_DOT    Section number requires trailing dot (e.g., '1.' not '1')
  SPACING         Exactly one space required after section number
  DEPTH_MISMATCH  Heading level doesn't match section number depth
  MISSING_PARENT  Parent section not defined before child
  ORDER           Section numbers not in ascending order
`, Version)
}

func printVersion() {
	fmt.Printf("md-section-numbers-checker %s\n", Version)
	fmt.Printf("  commit: %s\n", GitCommit)
	fmt.Printf("  built:  %s\n", BuildDate)
}

func processPatterns(patterns []string) bool {
	hasError := false
	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid pattern: %s\n", pattern)
			hasError = true
			continue
		}

		if len(files) == 0 {
			fmt.Fprintf(os.Stderr, "No files matched: %s\n", pattern)
			continue
		}

		for _, file := range files {
			if processFile(file) {
				hasError = true
			}
		}
	}
	return hasError
}

func processFile(file string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", file, err)
		return true
	}

	errors := checker.ValidateContent(string(content))
	for _, e := range errors {
		fmt.Fprintf(os.Stderr, "%s:%d: [%s] %s\n", file, e.Line, e.Code, e.Message)
	}

	return len(errors) > 0
}
