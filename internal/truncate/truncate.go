// Package truncate provides field-level value truncation for structured log lines.
// Long string values can be capped to a maximum byte length, with an optional
// ellipsis suffix to indicate the value was shortened.
package truncate

import (
	"encoding/json"
	"fmt"
)

// Truncator truncates string values in JSON log entries that exceed a maximum length.
type Truncator struct {
	maxLen int
	suffix string
}

// New creates a Truncator that caps string field values to maxLen bytes.
// If maxLen is zero or negative the Truncator is a no-op.
// suffix is appended when a value is shortened (e.g. "…").
func New(maxLen int, suffix string) (*Truncator, error) {
	if maxLen < 0 {
		return nil, fmt.Errorf("truncate: maxLen must be >= 0, got %d", maxLen)
	}
	return &Truncator{maxLen: maxLen, suffix: suffix}, nil
}

// Apply returns a copy of msg with all string values longer than maxLen truncated.
// Non-string values and the structure of the object are preserved.
// If maxLen is 0, msg is returned unchanged.
func (t *Truncator) Apply(msg json.RawMessage) (json.RawMessage, error) {
	if t.maxLen == 0 {
		return msg, nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return msg, nil // not a JSON object — pass through
	}

	modified := false
	for k, v := range obj {
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			continue // not a string field
		}
		if len(s) > t.maxLen {
			truncated := s[:t.maxLen] + t.suffix
			encoded, err := json.Marshal(truncated)
			if err != nil {
				return msg, fmt.Errorf("truncate: marshal field %q: %w", k, err)
			}
			obj[k] = json.RawMessage(encoded)
			modified = true
		}
	}

	if !modified {
		return msg, nil
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return msg, fmt.Errorf("truncate: marshal result: %w", err)
	}
	return out, nil
}
