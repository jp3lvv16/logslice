// Package coalesce provides a processor that merges a list of candidate
// fields into a single target field, using the first non-empty value found.
package coalesce

import (
	"encoding/json"
	"fmt"
)

// Coalescer merges candidate fields into a target field.
type Coalescer struct {
	target     string
	candidates []string
	keepSrc    bool
}

// New returns a Coalescer that reads from candidates (in order) and writes
// the first non-empty string value to target. If keepSrc is false the
// candidate fields that were consumed are removed from the object.
func New(target string, candidates []string, keepSrc bool) (*Coalescer, error) {
	if target == "" {
		return nil, fmt.Errorf("coalesce: target field name must not be empty")
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("coalesce: at least one candidate field is required")
	}
	return &Coalescer{target: target, candidates: candidates, keepSrc: keepSrc}, nil
}

// Apply reads raw JSON, coalesces the candidate fields into the target and
// returns the modified JSON. If no candidate contains a usable value the
// message is returned unchanged.
func (c *Coalescer) Apply(raw json.RawMessage) (json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return raw, nil
	}

	var chosen json.RawMessage
	var usedKey string

	for _, key := range c.candidates {
		v, ok := obj[key]
		if !ok {
			continue
		}
		// Reject JSON null and empty string.
		if string(v) == "null" || string(v) == `""` {
			continue
		}
		chosen = v
		usedKey = key
		break
	}

	if chosen == nil {
		return raw, nil
	}

	obj[c.target] = chosen

	if !c.keepSrc && usedKey != c.target {
		delete(obj, usedKey)
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return raw, err
	}
	return out, nil
}
