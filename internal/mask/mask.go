// Package mask provides field-level value masking using regular expressions.
// Matched portions of a field's string value are replaced with a fixed mask string.
package mask

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Rule pairs a compiled pattern with the mask string to substitute.
type Rule struct {
	field   string
	pattern *regexp.Regexp
	mask    string
}

// Masker holds a set of masking rules applied to JSON log lines.
type Masker struct {
	rules []Rule
}

// New parses specs of the form "field/pattern/mask" and returns a Masker.
// The separator is the first character after the field name's slash.
func New(specs []string) (*Masker, error) {
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, "/", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("mask: invalid spec %q: want field/pattern/mask", spec)
		}
		re, err := regexp.Compile(parts[1])
		if err != nil {
			return nil, fmt.Errorf("mask: invalid pattern in spec %q: %w", spec, err)
		}
		rules = append(rules, Rule{field: parts[0], pattern: re, mask: parts[2]})
	}
	return &Masker{rules: rules}, nil
}

// Apply returns a new JSON message with matching field values masked.
// Lines that are not valid JSON objects are returned unchanged.
func (m *Masker) Apply(line []byte) []byte {
	if len(m.rules) == 0 {
		return line
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(line, &obj); err != nil {
		return line
	}
	modified := false
	for _, r := range m.rules {
		raw, ok := obj[r.field]
		if !ok {
			continue
		}
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			continue
		}
		masked := r.pattern.ReplaceAllString(s, r.mask)
		if masked == s {
			continue
		}
		encoded, err := json.Marshal(masked)
		if err != nil {
			continue
		}
		obj[r.field] = json.RawMessage(encoded)
		modified = true
	}
	if !modified {
		return line
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return out
}
