// Package redact provides field-level value redaction for structured log lines.
// It replaces the value of nominated JSON fields with a fixed placeholder so
// that sensitive data (tokens, passwords, PII) is never written to output.
package redact

import (
	"encoding/json"
	"fmt"
	"strings"
)

const defaultPlaceholder = "[REDACTED]"

// Redactor replaces the values of a fixed set of top-level JSON fields with a
// placeholder string.
type Redactor struct {
	fields      map[string]struct{}
	placeholder string
}

// New returns a Redactor that will replace the values of the given field names.
// An optional placeholder may be supplied as the second argument; when omitted
// the default "[REDACTED]" is used.
func New(fields []string, placeholder string) (*Redactor, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("redact: at least one field name is required")
	}
	ph := strings.TrimSpace(placeholder)
	if ph == "" {
		ph = defaultPlaceholder
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			return nil, fmt.Errorf("redact: field name must not be empty")
		}
		set[f] = struct{}{}
	}
	return &Redactor{fields: set, placeholder: ph}, nil
}

// Apply redacts the configured fields in the raw JSON message and returns the
// modified bytes. Lines that are not valid JSON objects are returned unchanged.
func (r *Redactor) Apply(raw json.RawMessage) (json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		// Not a JSON object — pass through untouched.
		return raw, nil
	}
	for field := range r.fields {
		if _, ok := obj[field]; ok {
			quoted, _ := json.Marshal(r.placeholder)
			obj[field] = quoted
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return raw, fmt.Errorf("redact: re-marshal failed: %w", err)
	}
	return out, nil
}
