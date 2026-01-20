package checker

import (
	"testing"
)

// helper
func checkHeadingEqual(t *testing.T, result, expected *Heading) {
	if result.Level != expected.Level {
		t.Errorf("expected Level %v, got %v", expected.Level, result.Level)
	}
	if result.Numbering != expected.Numbering {
		t.Errorf("expected Numbering %v, got %v", expected.Numbering, result.Numbering)
	}
	if len(result.Segments) != len(expected.Segments) {
		t.Errorf("expected Segments length %v, got %v", expected.Segments, result.Segments)
	} else {
		for i := range result.Segments {
			if result.Segments[i] != expected.Segments[i] {
				t.Errorf("expected Segment %v, got %v", expected.Segments, result.Segments)
			}
		}
	}
	if result.RawNumbering != expected.RawNumbering {
		t.Errorf("expected RawNumbering %v, got %v", expected.RawNumbering, result.RawNumbering)
	}
	if result.Spacing != expected.Spacing {
		t.Errorf("expected Spacing %v, got %v", expected.Spacing, result.Spacing)
	}
}

// tests
func TestExtractHeading_Level2(t *testing.T) {
	line := "## 1.2.3. Title"
	expected := &Heading{
		Level:        2,
		Numbering:    "1.2.3",
		Segments:     []string{"1", "2", "3"},
		RawNumbering: "1.2.3.",
		Spacing:      " ",
	}
	result := extractHeading(line)
	if result == nil {
		t.Fatalf("expected heading, got nil")
	}
	checkHeadingEqual(t, result, expected)
}

func TestExtractHeading_Level6(t *testing.T) {
	line := "###### 1.2.3. Title"
	expected := &Heading{
		Level:        6,
		Numbering:    "1.2.3",
		Segments:     []string{"1", "2", "3"},
		RawNumbering: "1.2.3.",
		Spacing:      " ",
	}
	result := extractHeading(line)
	if result == nil {
		t.Fatalf("expected heading, got nil")
	}
	checkHeadingEqual(t, result, expected)
}

func TestExtractHeading_TooManyHashes(t *testing.T) {
	line := "####### 1.2.3. Title"

	result := extractHeading(line)
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}

func TestExtractHeading_NoSpace(t *testing.T) {
	line := "### 1.2.3.Title"
	expected := &Heading{
		Level:        3,
		Numbering:    "1.2.3",
		Segments:     []string{"1", "2", "3"},
		RawNumbering: "1.2.3.",
		Spacing:      "",
	}
	result := extractHeading(line)
	if result == nil {
		t.Fatalf("expected heading, got nil")
	}
	checkHeadingEqual(t, result, expected)
}

func TestExtractHeading_NoTitle(t *testing.T) {
	line := "### 1.2.3."
	expected := &Heading{
		Level:        3,
		Numbering:    "1.2.3",
		Segments:     []string{"1", "2", "3"},
		RawNumbering: "1.2.3.",
		Spacing:      "",
	}
	result := extractHeading(line)
	if result == nil {
		t.Fatalf("expected heading, got nil")
	}
	checkHeadingEqual(t, result, expected)
}

func TestExtractHeading_NotHeading(t *testing.T) {
	line := "1.2.3. Title"
	result := extractHeading(line)
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}

func TestExtractHeading_OnlyHashes(t *testing.T) {
	line := "#######"
	result := extractHeading(line)
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}

func TestExtractHeading_HasNoNumber(t *testing.T) {
	line := "## Title"
	result := extractHeading(line)
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}
func TestExtractHeading_NotNumber(t *testing.T) {
	line := "### a.b.c. Title"
	result := extractHeading(line)
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}

func TestExtractHeading_IncludeNotNumber(t *testing.T) {
	line := "### 1.2.a. Title"
	expected := &Heading{
		Level:        3,
		Numbering:    "1.2",
		Segments:     []string{"1", "2"},
		RawNumbering: "1.2.",
		Spacing:      "",
	}
	result := extractHeading(line)
	if result == nil {
		t.Fatalf("expected heading, got nil")
	}
	checkHeadingEqual(t, result, expected)
}
