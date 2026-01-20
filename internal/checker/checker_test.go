package checker

import (
	"strings"
	"testing"
)

func TestValidateContent_ValidDoc(t *testing.T) {
	validDoc := "# Title\n\n## TOC\n\n## 1. Parent\n\n### 1.1. Child\ncontent\n\n### 1.2. Child\ncontent\n\n## 2. Next Parent\n"
	errors := ValidateContent(validDoc)
	if len(errors) != 0 {
		t.Errorf("expected no errors, got %d: %v", len(errors), errors)
	}
}

func TestValidateContent_LevelDepthMismatch(t *testing.T) {
	invalidDoc := "# Title\n\n### 1. Child\n"
	errors := ValidateContent(invalidDoc)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Code != CodeDepthMismatch {
		t.Errorf("expected code %q, got %q", CodeDepthMismatch, errors[0].Code)
	}
	if !strings.Contains(errors[0].Message, "does not match section number depth") {
		t.Errorf("expected message to contain 'does not match section number depth', got %q", errors[0].Message)
	}
	if errors[0].Line != 3 {
		t.Errorf("expected line 3, got %d", errors[0].Line)
	}
}

func TestValidateContent_ParentNotDefinedBeforeChild(t *testing.T) {
	invalidDoc := "# Title\n\n### 1.1. Child\n\n## 1. Parent\n"
	errors := ValidateContent(invalidDoc)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Code != CodeMissingParent {
		t.Errorf("expected code %q, got %q", CodeMissingParent, errors[0].Code)
	}
	if !strings.Contains(errors[0].Message, "before parent section 1 is defined") {
		t.Errorf("expected message to contain 'before parent section 1 is defined', got %q", errors[0].Message)
	}
}

func TestValidateContent_NotAscending(t *testing.T) {
	invalidDoc := "# Title\n\n## 2. Second\n\n## 1. First\n"
	errors := ValidateContent(invalidDoc)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Code != CodeOrder {
		t.Errorf("expected code %q, got %q", CodeOrder, errors[0].Code)
	}
	if !strings.Contains(errors[0].Message, "top-level") {
		t.Errorf("expected message to contain 'top-level', got %q", errors[0].Message)
	}
}

func TestValidateContent_MissingTrailingDot(t *testing.T) {
	invalidDoc := "# Title\n\n## 1 Parent\n"
	errors := ValidateContent(invalidDoc)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Code != CodeTrailingDot {
		t.Errorf("expected code %q, got %q", CodeTrailingDot, errors[0].Code)
	}
	if !strings.Contains(errors[0].Message, "requires a trailing dot") {
		t.Errorf("expected message to contain 'requires a trailing dot', got %q", errors[0].Message)
	}
}

func TestValidateContent_MissingSpaceAfterNumber(t *testing.T) {
	invalidDoc := "# Title\n\n## 1.Parent\n"
	errors := ValidateContent(invalidDoc)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Code != CodeSpacing {
		t.Errorf("expected code %q, got %q", CodeSpacing, errors[0].Code)
	}
	if !strings.Contains(errors[0].Message, "must be followed by exactly one space") {
		t.Errorf("expected message to contain 'must be followed by exactly one space', got %q", errors[0].Message)
	}
}

func TestValidateContent_TooManySpacesAfterNumber(t *testing.T) {
	invalidDoc := "# Title\n\n## 1.  Parent\n"
	errors := ValidateContent(invalidDoc)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
	if errors[0].Code != CodeSpacing {
		t.Errorf("expected code %q, got %q", CodeSpacing, errors[0].Code)
	}
	if !strings.Contains(errors[0].Message, "must be followed by exactly one space") {
		t.Errorf("expected message to contain 'must be followed by exactly one space', got %q", errors[0].Message)
	}
}
