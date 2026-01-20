package checker

import (
	"regexp"
	"strings"
)

// Heading represents a parsed Markdown heading with section number.
type Heading struct {
	Level        int
	Numbering    string
	Segments     []string
	RawNumbering string
	Spacing      string
}

// headingRegex matches Markdown headings with section numbers.
//
// Pattern: ^(#{2,6})\s+(\d+(?:\.\d+)*\.?)(\s*)
//   - (#{2,6})           : h2 to h6 (## to ######)
//   - \s+                : one or more whitespace
//   - (\d+(?:\.\d+)*\.?) : section number (e.g., "1", "1.2", "1.2.3.")
//   - (\s*)              : spacing after the number (captured for validation)
var headingRegex = regexp.MustCompile(`^(#{2,6})\s+(\d+(?:\.\d+)*\.?)(\s*)`)

// extractHeading parses a line and returns a Heading if it matches.
//
// Example: "### 1.2. Title" returns:
//
//	&Heading{
//	    Level:        3,          // ### = h3
//	    Numbering:    "1.2",      // without trailing dot
//	    Segments:     ["1", "2"], // split by dot
//	    RawNumbering: "1.2.",     // original string
//	    Spacing:      " ",        // space after number
//	}
//
// Returns nil if the line does not match (e.g., h1, no number, non-numeric).
func extractHeading(line string) *Heading {
	matches := headingRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	level := len(matches[1])
	rawNumbering := matches[2]
	spacing := matches[3]
	numbering := strings.TrimSuffix(rawNumbering, ".")
	segments := strings.Split(numbering, ".")

	// Filter empty segments
	filtered := make([]string, 0, len(segments))
	for _, s := range segments {
		if s != "" {
			filtered = append(filtered, s)
		}
	}

	return &Heading{
		Level:        level,
		Numbering:    numbering,
		Segments:     filtered,
		RawNumbering: rawNumbering,
		Spacing:      spacing,
	}
}
