// Package normalize rewrites field values to a canonical form.
// It supports case folding (lower, upper, title) and whitespace
// trimming for string fields in a JSON log message.
package normalize

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Mode controls how a string value is normalized.
type Mode string

const (
	ModeLower Mode = "lower"
	ModeUpper Mode = "upper"
	ModeTitle Mode = "title"
	ModeTrim  Mode = "trim"
)

// rule pairs a field name with the normalization mode to apply.
type rule struct {
	field string
	mode  Mode
}

// Normalizer applies normalization rules to JSON log messages.
type Normalizer struct {
	rules []rule
}

// New constructs a Normalizer from a slice of specs.
// Each spec must be in the form "field:mode" where mode is one of
// lower, upper, title, or trim.
func New(specs []string) (*Normalizer, error) {
	if len(specs) == 0 {
		return &Normalizer{}, nil
	}
	rules := make([]rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("normalize: invalid spec %q: expected field:mode", s)
		}
		field := strings.TrimSpace(parts[0])
		mode := Mode(strings.TrimSpace(parts[1]))
		if field == "" {
			return nil, fmt.Errorf("normalize: empty field name in spec %q", s)
		}
		switch mode {
		case ModeLower, ModeUpper, ModeTitle, ModeTrim:
		default:
			return nil, fmt.Errorf("normalize: unknown mode %q in spec %q", mode, s)
		}
		rules = append(rules, rule{field: field, mode: mode})
	}
	return &Normalizer{rules: rules}, nil
}

// Apply rewrites the configured fields in msg and returns the updated JSON.
// Fields that are absent or non-string are left untouched.
func (n *Normalizer) Apply(msg json.RawMessage) (json.RawMessage, error) {
	if len(n.rules) == 0 {
		return msg, nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return msg, nil
	}
	for _, r := range n.rules {
		raw, ok := obj[r.field]
		if !ok {
			continue
		}
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			continue
		}
		switch r.mode {
		case ModeLower:
			s = strings.ToLower(s)
		case ModeUpper:
			s = strings.ToUpper(s)
		case ModeTitle:
			s = strings.ToTitle(s)
		case ModeTrim:
			s = strings.TrimSpace(s)
		}
		b, _ := json.Marshal(s)
		obj[r.field] = b
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return msg, err
	}
	return out, nil
}
