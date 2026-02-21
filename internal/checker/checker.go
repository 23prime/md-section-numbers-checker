package checker

import (
	"fmt"
	"strconv"
	"strings"
)

// Error codes
const (
	CodeTrailingDot   = "TRAILING_DOT"
	CodeSpacing       = "SPACING"
	CodeDepthMismatch = "DEPTH_MISMATCH"
	CodeMissingParent = "MISSING_PARENT"
	CodeOrder         = "ORDER"
)

// rootKey is the key used for top-level section ordering.
const rootKey = "root"

// Error represents a validation error.
type Error struct {
	Line    int
	Code    string
	Message string
}

// NewError creates a new validation error.
func NewError(line int, code string, message string) Error {
	return Error{
		Line:    line,
		Code:    code,
		Message: message,
	}
}

// ValidateContent validates Markdown content and returns errors.
func ValidateContent(content string) []Error {
	lines := strings.Split(content, "\n")
	// Handle CRLF
	for i, line := range lines {
		lines[i] = strings.TrimSuffix(line, "\r")
	}
	return validateLines(lines)
}

func validateLines(lines []string) []Error {
	var errors []Error
	state := newValidationState()

	for idx, line := range lines {
		lineNumber := idx + 1
		heading := extractHeading(line)
		if heading == nil {
			continue
		}

		headingErrors := validateHeading(state, heading, lineNumber)
		errors = append(errors, headingErrors...)
	}

	return errors
}

func validateHeading(state *validationState, heading *Heading, lineNumber int) []Error {
	var errors []Error

	// Check trailing dot
	if err := checkTrailingDot(heading, lineNumber); err != nil {
		errors = append(errors, *err)
	}

	// Check spacing
	if err := checkSpacing(heading, lineNumber); err != nil {
		errors = append(errors, *err)
	}

	// Check level-depth consistency
	if err := checkDepthConsistency(heading, lineNumber); err != nil {
		errors = append(errors, *err)
		return errors
	}

	// Compute parent info
	parent := computeParentInfo(heading)

	// Check parent exists
	if err := checkParentExists(state, heading, parent, lineNumber); err != nil {
		errors = append(errors, *err)
		return errors
	}

	// Parse last segment
	currentValue, err := parseLastSegment(heading, lineNumber)
	if err != nil {
		errors = append(errors, *err)
		return errors
	}

	// Check ascending order
	if err := checkAscendingOrder(state, parent, currentValue, lineNumber); err != nil {
		errors = append(errors, *err)
		return errors
	}

	// Update state
	updateState(state, heading, parent, currentValue)

	return errors
}

// validationState holds the state during validation.
//
// Fields:
//   - seenHeadings: tracks which section numbers have been defined.
//     Used to check if a parent section exists before its child.
//     Example: {"1": true, "1.1": true, "1.2": true, "2": true}
//   - lastCounters: tracks the last section number at each depth level.
//     Used to verify ascending order within the same parent.
//     Key is parent section (or "root" for top-level).
//     Example: {"root": 2, "1": 2} means top-level is at 2, and under "1" is at 1.2.
type validationState struct {
	seenHeadings map[string]bool
	lastCounters map[string]int
}

func newValidationState() *validationState {
	return &validationState{
		seenHeadings: make(map[string]bool),
		lastCounters: make(map[string]int),
	}
}

func checkTrailingDot(heading *Heading, lineNumber int) *Error {
	if !strings.HasSuffix(heading.RawNumbering, ".") {
		err := NewError(
			lineNumber,
			CodeTrailingDot,
			fmt.Sprintf("section number %s requires a trailing dot (e.g., %s.)", heading.RawNumbering, heading.Numbering),
		)
		return &err
	}
	return nil
}

func checkSpacing(heading *Heading, lineNumber int) *Error {
	if heading.Spacing != " " {
		err := NewError(
			lineNumber,
			CodeSpacing,
			fmt.Sprintf("section number %s must be followed by exactly one space", heading.RawNumbering),
		)
		return &err
	}
	return nil
}

func checkDepthConsistency(heading *Heading, lineNumber int) *Error {
	expectedSegments := heading.Level - 1
	if len(heading.Segments) != expectedSegments {
		err := NewError(
			lineNumber,
			CodeDepthMismatch,
			fmt.Sprintf("heading level (h%d) does not match section number depth %s", heading.Level, heading.Numbering),
		)
		return &err
	}
	return nil
}

// parentInfo holds computed parent section information.
type parentInfo struct {
	segments []string
	key      string
	depthKey string
}

func computeParentInfo(heading *Heading) parentInfo {
	segments := heading.Segments[:len(heading.Segments)-1]
	key := strings.Join(segments, ".")
	depthKey := rootKey
	if len(segments) > 0 {
		depthKey = key
	}
	return parentInfo{
		segments: segments,
		key:      key,
		depthKey: depthKey,
	}
}

func checkParentExists(state *validationState, heading *Heading, parent parentInfo, lineNumber int) *Error {
	if len(parent.segments) > 0 && !state.seenHeadings[parent.key] {
		err := NewError(
			lineNumber,
			CodeMissingParent,
			fmt.Sprintf("child section %s appears before parent section %s is defined", heading.Numbering, parent.key),
		)
		return &err
	}
	return nil
}

func parseLastSegment(heading *Heading, lineNumber int) (int, *Error) {
	lastSegment := heading.Segments[len(heading.Segments)-1]
	value, err := strconv.Atoi(lastSegment)
	if err != nil {
		e := NewError(
			lineNumber,
			CodeOrder,
			fmt.Sprintf("section number %s has non-numeric segment", heading.Numbering),
		)
		return 0, &e
	}
	return value, nil
}

func checkAscendingOrder(state *validationState, parent parentInfo, currentValue int, lineNumber int) *Error {
	scope := "at top-level"
	if len(parent.segments) > 0 {
		scope = fmt.Sprintf("under parent section %s", parent.key)
	}
	previousValue, exists := state.lastCounters[parent.depthKey]
	if !exists {
		if currentValue != 1 {
			e := NewError(
				lineNumber,
				CodeOrder,
				fmt.Sprintf("headings %s must start at 1 (got: %d)", scope, currentValue),
			)
			return &e
		}
		return nil
	}
	if currentValue != previousValue+1 {
		e := NewError(
			lineNumber,
			CodeOrder,
			fmt.Sprintf("headings %s are not consecutive (expected: %d, got: %d)", scope, previousValue+1, currentValue),
		)
		return &e
	}
	return nil
}

func updateState(state *validationState, heading *Heading, parent parentInfo, currentValue int) {
	state.lastCounters[parent.depthKey] = currentValue
	state.seenHeadings[strings.Join(heading.Segments, ".")] = true
}
