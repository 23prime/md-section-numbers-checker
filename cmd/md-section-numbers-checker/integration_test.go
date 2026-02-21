package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "md-section-numbers-checker-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmp)

	binPath = filepath.Join(tmp, "md-section-numbers-checker")
	if err := exec.Command("go", "build", "-o", binPath, ".").Run(); err != nil {
		panic("failed to build binary: " + err.Error())
	}

	os.Exit(m.Run())
}

func run(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(binPath, args...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	return
}

func TestHelp(t *testing.T) {
	for _, flag := range []string{"--help", "-h"} {
		t.Run(flag, func(t *testing.T) {
			stdout, _, exitCode := run(t, flag)
			if exitCode != 0 {
				t.Errorf("expected exit code 0, got %d", exitCode)
			}
			if !strings.Contains(stdout, "USAGE:") {
				t.Errorf("expected USAGE: in stdout, got: %s", stdout)
			}
			if !strings.Contains(stdout, "ERROR CODES:") {
				t.Errorf("expected ERROR CODES: in stdout, got: %s", stdout)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	stdout, _, exitCode := run(t, "--version")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "md-section-numbers-checker") {
		t.Errorf("expected binary name in stdout, got: %s", stdout)
	}
}

func TestNoArgs(t *testing.T) {
	_, stderr, exitCode := run(t)
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "Usage:") {
		t.Errorf("expected usage message in stderr, got: %s", stderr)
	}
}

func TestValidFile(t *testing.T) {
	f, err := os.CreateTemp("", "valid-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString("# Title\n\n## 1. Introduction\n\n### 1.1. Background\n\n## 2. Conclusion\n"); err != nil {
		t.Fatal(err)
	}
	f.Close()

	stdout, stderr, exitCode := run(t, f.Name())
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nstderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "All section numbers are valid.") {
		t.Errorf("expected success message in stdout, got: %s", stdout)
	}
}

func TestInvalidFile(t *testing.T) {
	f, err := os.CreateTemp("", "invalid-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	content := "# Title\n\n" +
		"## 1 Missing Trailing Dot\n\n" +
		"## 2.No Space After Dot\n\n" +
		"### 2.2. Not start at \"1\"\n\n" +
		"## 3.  Two Spaces After Dot\n\n" +
		"### 4. Depth Mismatch\n\n" +
		"### 5.1. Missing Parent\n\n" +
		"## 6. Not Consecutive After Gap\n\n" +
		"### 6.1. First Child\n\n" +
		"### 6.1. Not Ascending\n\n" +
		"## 8. Not Consecutive\n"
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	_, stderr, exitCode := run(t, f.Name())
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	for _, code := range []string{"TRAILING_DOT", "SPACING", "DEPTH_MISMATCH", "MISSING_PARENT", "ORDER"} {
		if !strings.Contains(stderr, code) {
			t.Errorf("expected error code %s in stderr, got: %s", code, stderr)
		}
	}
}

func TestNoFilesMatched(t *testing.T) {
	_, stderr, exitCode := run(t, "nonexistent.md")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stderr, "No files matched") {
		t.Errorf("expected 'No files matched' in stderr, got: %s", stderr)
	}
}
