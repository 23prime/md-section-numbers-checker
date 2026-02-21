package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		fmt.Fprintln(os.Stderr, "Usage: mdsnc <file.md> ...")
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
	fmt.Printf(`mdsnc %s

USAGE:
  mdsnc [OPTIONS] <file.md> ...

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
	fmt.Printf("mdsnc %s\n", Version)
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

		for _, file := range filterGitIgnored(files) {
			if processFile(file) {
				hasError = true
			}
		}
	}
	return hasError
}

// filterGitIgnored returns only the files that are NOT git-ignored.
// If git is unavailable or the directory is not a git repo, returns files unchanged.
func filterGitIgnored(files []string) []string {
	if len(files) == 0 {
		return files
	}
	if _, err := exec.LookPath("git"); err != nil {
		return files
	}

	var input strings.Builder
	for _, f := range files {
		input.WriteString(f)
		input.WriteByte('\n')
	}

	cmd := exec.Command("git", "check-ignore", "--stdin")
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// exit 1: no files ignored; exit 128: not a git repo — process all files
		return files
	}

	ignored := make(map[string]bool)
	for _, line := range strings.Split(out.String(), "\n") {
		if line != "" {
			ignored[line] = true
		}
	}

	result := make([]string, 0, len(files))
	for _, f := range files {
		if !ignored[f] {
			result = append(result, f)
		}
	}
	return result
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
