// Package transform provides field renaming and value rewriting for log records.
package transform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule describes a single transformation: rename a field or rewrite its value.
type Rule struct {
	FromKey string
	ToKey   string
	FromVal string // empty means match any value
	ToVal   string
}

// Transformer applies a set of Rules to raw JSON log lines.
type Transformer struct {
	rules []Rule
}

// New parses specs of the form "oldKey=newKey" (rename) or
// "key:oldVal=newVal" (value rewrite) and returns a Transformer.
func New(specs []string) (*Transformer, error) {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		r, err := parseSpec(s)
		if err != nil {
			return nil, fmt.Errorf("transform: %w", err)
		}
		rules = append(rules, r)
	}
	return &Transformer{rules: rules}, nil
}

// Apply runs all rules against the raw JSON message and returns the result.
// Lines that are not valid JSON objects are returned unchanged.
func (t *Transformer) Apply(raw json.RawMessage) json.RawMessage {
	if len(t.rules) == 0 {
		return raw
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return raw
	}
	for _, r := range t.rules {
		applyRule(m, r)
	}
	out, err := json.Marshal(m)
	if err != nil {
		return raw
	}
	return out
}

func applyRule(m map[string]json.RawMessage, r Rule) {
	val, ok := m[r.FromKey]
	if !ok {
		return
	}
	if r.FromVal != "" {
		// value rewrite — key stays the same
		var s string
		if err := json.Unmarshal(val, &s); err != nil || s != r.FromVal {
			return
		}
		newVal, _ := json.Marshal(r.ToVal)
		m[r.FromKey] = newVal
		return
	}
	// rename
	if r.FromKey != r.ToKey {
		m[r.ToKey] = val
		delete(m, r.FromKey)
	}
}

// parseSpec handles "oldKey=newKey" and "key:oldVal=newVal".
func parseSpec(s string) (Rule, error) {
	if colonIdx := strings.Index(s, ":"); colonIdx != -1 {
		key := s[:colonIdx]
		rest := s[colonIdx+1:]
		parts := strings.SplitN(rest, "=", 2)
		if len(parts) != 2 {
			return Rule{}, fmt.Errorf("invalid value-rewrite spec %q: expected key:oldVal=newVal", s)
		}
		return Rule{FromKey: key, ToKey: key, FromVal: parts[0], ToVal: parts[1]}, nil
	}
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("invalid rename spec %q: expected oldKey=newKey", s)
	}
	return Rule{FromKey: parts[0], ToKey: parts[1]}, nil
}
