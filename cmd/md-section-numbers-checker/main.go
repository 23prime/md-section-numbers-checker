package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/23prime/md-section-numbers-checker/internal/checker"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "show version information")
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: md-section-numbers-checker <file.md> ...")
		os.Exit(1)
	}

	hasError := processPatterns(args)

	if hasError {
		os.Exit(1)
	}

	fmt.Println("All section numbers are valid.")
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
