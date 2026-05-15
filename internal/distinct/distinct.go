// Package distinct provides a filter that drops duplicate values of a
// specified JSON field, emitting only the first occurrence of each unique value
// within an optional sliding window.
package distinct

import (
	"encoding/json"
	"sync"
)

// Distinct tracks unique values for a given field key and reports whether
// a message should be kept (i.e. its field value has not been seen before).
type Distinct struct {
	mu      sync.Mutex
	field   string
	seen    map[string]struct{}
	window  int // max unique values to remember; 0 = unlimited
	order   []string
}

// New creates a Distinct filter that tracks unique values of field.
// window controls how many unique values are remembered (0 = unlimited).
// Returns an error if field is empty.
func New(field string, window int) (*Distinct, error) {
	if field == "" {
		return nil, errEmptyField
	}
	if window < 0 {
		window = 0
	}
	return &Distinct{
		field:  field,
		seen:   make(map[string]struct{}),
		window: window,
	}, nil
}

// Keep returns true if the value of the configured field in msg has not been
// seen before. Malformed JSON or a missing field always returns true (pass-through).
func (d *Distinct) Keep(msg json.RawMessage) bool {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return true
	}
	raw, ok := obj[d.field]
	if !ok {
		return true
	}
	key := string(raw)

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.seen[key]; exists {
		return false
	}
	d.seen[key] = struct{}{}
	d.order = append(d.order, key)

	if d.window > 0 && len(d.order) > d.window {
		evict := d.order[0]
		d.order = d.order[1:]
		delete(d.seen, evict)
	}
	return true
}

// Len returns the number of unique values currently tracked.
func (d *Distinct) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
