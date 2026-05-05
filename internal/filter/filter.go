package filter

import (
	"regexp"
	"strings"
)

// Rule represents a single field-pattern filter rule.
type Rule struct {
	Field   string
	Pattern *regexp.Regexp
}

// Filter holds a set of rules used to match log entries.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of "field=pattern" strings.
// Each pattern is treated as a regular expression.
func New(specs []string) (*Filter, error) {
	f := &Filter{}
	for _, spec := range specs {
		idx := strings.IndexByte(spec, '=')
		if idx < 1 {
			return nil, &ParseError{Spec: spec, Reason: "missing '=' separator"}
		}
		field := spec[:idx]
		patternStr := spec[idx+1:]
		re, err := regexp.Compile(patternStr)
		if err != nil {
			return nil, &ParseError{Spec: spec, Reason: "invalid regex: " + err.Error()}
		}
		f.rules = append(f.rules, Rule{Field: field, Pattern: re})
	}
	return f, nil
}

// Match reports whether the given log entry (as a map of fields) satisfies
// all rules in the filter. An entry with no rules always matches.
func (f *Filter) Match(entry map[string]string) bool {
	for _, rule := range f.rules {
		val, ok := entry[rule.Field]
		if !ok || !rule.Pattern.MatchString(val) {
			return false
		}
	}
	return true
}

// Rules returns a copy of the filter's rules. This is useful for inspection
// or debugging without exposing the internal slice directly.
func (f *Filter) Rules() []Rule {
	result := make([]Rule, len(f.rules))
	copy(result, f.rules)
	return result
}

// ParseError is returned when a filter spec cannot be parsed.
type ParseError struct {
	Spec   string
	Reason string
}

func (e *ParseError) Error() string {
	return "filter: invalid spec " + e.Spec + ": " + e.Reason
}
