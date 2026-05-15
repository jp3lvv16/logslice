// Package fieldmap renames a set of JSON fields according to a mapping spec.
// Each spec has the form "old:new"; the old key is removed and its value is
// written under the new key. If the old key is absent the message is passed
// through unchanged.
package fieldmap

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Mapper holds a set of rename rules derived from the caller-supplied specs.
type Mapper struct {
	rules map[string]string // old -> new
}

// New parses specs of the form "old:new" and returns a Mapper.
// An error is returned if any spec is malformed or if old == new.
func New(specs []string) (*Mapper, error) {
	rules := make(map[string]string, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("fieldmap: invalid spec %q: want \"old:new\"", s)
		}
		if parts[0] == parts[1] {
			return nil, fmt.Errorf("fieldmap: spec %q maps a field to itself", s)
		}
		rules[parts[0]] = parts[1]
	}
	return &Mapper{rules: rules}, nil
}

// Apply renames fields in msg according to the mapping rules and returns the
// resulting JSON object. Fields that are not mentioned in the rules are kept
// as-is. If msg is not a valid JSON object it is returned unmodified.
func (m *Mapper) Apply(msg json.RawMessage) json.RawMessage {
	if len(m.rules) == 0 {
		return msg
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return msg
	}

	for old, newKey := range m.rules {
		val, ok := obj[old]
		if !ok {
			continue
		}
		delete(obj, old)
		obj[newKey] = val
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return msg
	}
	return out
}
