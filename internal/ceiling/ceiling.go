// Package ceiling caps the number of times a given field value may appear
// in the output stream. Once a value has been seen [Max] times it is dropped
// from subsequent messages; messages whose target field is absent are always
// forwarded unchanged.
package ceiling

import (
	"encoding/json"
	"errors"
	"sync"
)

// Ceiling tracks per-value occurrence counts and drops messages once a
// configured maximum has been reached for a given field value.
type Ceiling struct {
	field string
	max   int
	mu    sync.Mutex
	counts map[string]int
}

// New creates a Ceiling that allows at most max occurrences of each distinct
// value of field. max must be greater than zero.
func New(field string, max int) (*Ceiling, error) {
	if field == "" {
		return nil, errors.New("ceiling: field must not be empty")
	}
	if max <= 0 {
		return nil, errors.New("ceiling: max must be greater than zero")
	}
	return &Ceiling{
		field:  field,
		max:    max,
		counts: make(map[string]int),
	}, nil
}

// Allow reports whether msg should be forwarded. It returns true when the
// value of the configured field in msg has been seen fewer than max times
// (or when the field is absent). Each call to Allow that returns true
// increments the internal counter for that value.
func (c *Ceiling) Allow(msg json.RawMessage) bool {
	var rec map[string]json.RawMessage
	if err := json.Unmarshal(msg, &rec); err != nil {
		return true
	}
	raw, ok := rec[c.field]
	if !ok {
		return true
	}
	key := string(raw)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.counts[key] >= c.max {
		return false
	}
	c.counts[key]++
	return true
}

// Seen returns how many times value has been observed for the configured field.
func (c *Ceiling) Seen(value string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[value]
}
