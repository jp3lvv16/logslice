// Package flatten provides utilities for flattening nested JSON log objects
// into a single-level map using dot-notation keys.
//
// Example:
//
//	{"a": {"b": 1}} → {"a.b": 1}
package flatten

import (
	"encoding/json"
	"fmt"
)

// Flattener flattens nested JSON messages.
type Flattener struct {
	separator string
}

// New returns a Flattener using the given key separator (e.g. ".").
func New(separator string) *Flattener {
	if separator == "" {
		separator = "."
	}
	return &Flattener{separator: separator}
}

// Apply takes a raw JSON message and returns a flattened JSON message.
// Non-object top-level values are returned unchanged.
func (f *Flattener) Apply(msg json.RawMessage) (json.RawMessage, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(msg, &top); err != nil {
		// Not a JSON object — return as-is.
		return msg, nil
	}

	out := make(map[string]interface{})
	f.flatten(top, "", out)

	result, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("flatten: marshal: %w", err)
	}
	return result, nil
}

// flatten recursively walks nested objects, building dot-separated keys.
func (f *Flattener) flatten(obj map[string]json.RawMessage, prefix string, out map[string]interface{}) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = prefix + f.separator + k
		}

		var nested map[string]json.RawMessage
		if err := json.Unmarshal(v, &nested); err == nil {
			f.flatten(nested, key, out)
			continue
		}

		var scalar interface{}
		if err := json.Unmarshal(v, &scalar); err == nil {
			out[key] = scalar
		}
	}
}
