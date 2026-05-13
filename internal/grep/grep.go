// Package grep provides full-text search across raw JSON log lines.
// Unlike field-level filters, grep matches against the entire serialised
// log line, making it useful for quick ad-hoc searches.
package grep

import (
	"fmt"
	"regexp"
)

// Grep holds a compiled regular expression used to match raw log lines.
type Grep struct {
	re      *regexp.Regexp
	invert  bool
}

// New compiles pattern into a Grep matcher.
// If invert is true, Keep returns true only when the pattern does NOT match.
func New(pattern string, invert bool) (*Grep, error) {
	if pattern == "" {
		return nil, fmt.Errorf("grep: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("grep: invalid pattern %q: %w", pattern, err)
	}
	return &Grep{re: re, invert: invert}, nil
}

// Keep reports whether raw should be kept.
// raw is expected to be the raw bytes of a single log line.
func (g *Grep) Keep(raw []byte) bool {
	matched := g.re.Match(raw)
	if g.invert {
		return !matched
	}
	return matched
}

// Pattern returns the source pattern string.
func (g *Grep) Pattern() string {
	return g.re.String()
}

// Inverted reports whether the matcher is operating in invert mode.
func (g *Grep) Inverted() bool {
	return g.invert
}
