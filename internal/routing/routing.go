// Package routing provides field-based log routing — directing log lines
// to named output buckets based on field value patterns.
package routing

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Rule maps a compiled regex against a field value to a named destination.
type Rule struct {
	field   string
	pattern *regexp.Regexp
	dest    string
}

// Router holds an ordered list of routing rules.
type Router struct {
	rules       []Rule
	defaultDest string
}

// New creates a Router from specs of the form "field=pattern:dest".
// defaultDest is used when no rule matches.
func New(specs []string, defaultDest string) (*Router, error) {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		r, err := parseSpec(s)
		if err != nil {
			return nil, fmt.Errorf("routing: invalid spec %q: %w", s, err)
		}
		rules = append(rules, r)
	}
	return &Router{rules: rules, defaultDest: defaultDest}, nil
}

// Route returns the destination name for the given JSON log line.
// It returns defaultDest if no rule matches or the line is not valid JSON.
func (r *Router) Route(line []byte) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(line, &obj); err != nil {
		return r.defaultDest
	}
	for _, rule := range r.rules {
		raw, ok := obj[rule.field]
		if !ok {
			continue
		}
		var val string
		if err := json.Unmarshal(raw, &val); err != nil {
			continue
		}
		if rule.pattern.MatchString(val) {
			return rule.dest
		}
	}
	return r.defaultDest
}

// parseSpec parses "field=pattern:dest".
func parseSpec(s string) (Rule, error) {
	// find last colon as dest separator
	colon := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' {
			colon = i
			break
		}
	}
	if colon < 1 {
		return Rule{}, fmt.Errorf("missing ':dest' suffix")
	}
	dest := s[colon+1:]
	if dest == "" {
		return Rule{}, fmt.Errorf("dest must not be empty")
	}
	lhs := s[:colon]
	eq := -1
	for i, c := range lhs {
		if c == '=' {
			eq = i
			break
		}
	}
	if eq < 1 {
		return Rule{}, fmt.Errorf("missing 'field=' prefix")
	}
	field := lhs[:eq]
	pattern := lhs[eq+1:]
	re, err := regexp.Compile(pattern)
	if err != nil {
		return Rule{}, fmt.Errorf("invalid regex %q: %w", pattern, err)
	}
	return Rule{field: field, pattern: re, dest: dest}, nil
}
