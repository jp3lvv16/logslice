// Package topk tracks the top-K most frequent values for a given JSON field.
package topk

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Entry holds a field value and its observed count.
type Entry struct {
	Value string
	Count int
}

// Tracker accumulates value frequencies and returns the top-K entries.
type Tracker struct {
	field string
	k     int
	counts map[string]int
}

// New creates a Tracker that watches field and keeps the top k entries.
// k must be >= 1 and field must be non-empty.
func New(field string, k int) (*Tracker, error) {
	if field == "" {
		return nil, fmt.Errorf("topk: field must not be empty")
	}
	if k < 1 {
		return nil, fmt.Errorf("topk: k must be >= 1, got %d", k)
	}
	return &Tracker{
		field:  field,
		k:      k,
		counts: make(map[string]int),
	}, nil
}

// Add records the value of the configured field from a raw JSON message.
// Malformed JSON and missing fields are silently skipped.
func (t *Tracker) Add(msg json.RawMessage) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return
	}
	raw, ok := obj[t.field]
	if !ok {
		return
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		// fall back to raw token representation
		s = string(raw)
	}
	t.counts[s]++
}

// Top returns up to k entries sorted by descending count.
func (t *Tracker) Top() []Entry {
	entries := make([]Entry, 0, len(t.counts))
	for v, c := range t.counts {
		entries = append(entries, Entry{Value: v, Count: c})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	if len(entries) > t.k {
		entries = entries[:t.k]
	}
	return entries
}

// Reset clears all accumulated counts.
func (t *Tracker) Reset() {
	t.counts = make(map[string]int)
}
