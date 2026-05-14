// Package labelset provides a processor that attaches a fixed set of
// key/value labels to every log message that passes through it.
// Labels are specified as "key=value" pairs and are only applied when
// the target key is not already present in the message (non-overwriting).
package labelset

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Labeler attaches static labels to JSON log messages.
type Labeler struct {
	labels map[string]string
}

// New parses a slice of "key=value" spec strings and returns a Labeler.
// It returns an error if any spec is missing the "=" separator or has an
// empty key.
func New(specs []string) (*Labeler, error) {
	labels := make(map[string]string, len(specs))
	for _, s := range specs {
		key, val, ok := strings.Cut(s, "=")
		if !ok {
			return nil, fmt.Errorf("labelset: missing '=' in spec %q", s)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("labelset: empty key in spec %q", s)
		}
		labels[key] = val
	}
	return &Labeler{labels: labels}, nil
}

// Apply attaches configured labels to msg. Existing keys are never
// overwritten. The returned slice is a new JSON object. If msg is not
// valid JSON the original bytes are returned unchanged together with the
// parse error.
func (l *Labeler) Apply(msg json.RawMessage) (json.RawMessage, error) {
	if len(l.labels) == 0 {
		return msg, nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return msg, err
	}

	for k, v := range l.labels {
		if _, exists := obj[k]; !exists {
			quoted, _ := json.Marshal(v)
			obj[k] = quoted
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return msg, err
	}
	return out, nil
}
