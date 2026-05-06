// Package aggregate counts occurrences of a field value across log lines,
// enabling simple frequency analysis from the terminal.
package aggregate

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// Counter accumulates counts keyed by a JSON field value.
type Counter struct {
	field  string
	counts map[string]int
}

// New returns a Counter that tracks occurrences of the given field.
// An empty field name is valid; Add will simply record every line under
// the empty-string bucket.
func New(field string) *Counter {
	return &Counter{
		field:  field,
		counts: make(map[string]int),
	}
}

// Add extracts the target field from raw JSON and increments its counter.
// Lines that are not valid JSON or that lack the field are counted under
// the special key "<missing>".
func (c *Counter) Add(raw json.RawMessage) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		c.counts["<missing>"]++
		return
	}
	v, ok := obj[c.field]
	if !ok {
		c.counts["<missing>"]++
		return
	}
	// Unquote strings; keep other types as their JSON representation.
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		s = string(v)
	}
	c.counts[s]++
}

// Entry is a single aggregation result.
type Entry struct {
	Value string
	Count int
}

// Results returns entries sorted by count descending, then value ascending.
func (c *Counter) Results() []Entry {
	entries := make([]Entry, 0, len(c.counts))
	for k, v := range c.counts {
		entries = append(entries, Entry{Value: k, Count: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	return entries
}

// Print writes a human-readable frequency table to w.
func (c *Counter) Print(w io.Writer) {
	for _, e := range c.Results() {
		fmt.Fprintf(w, "%6d  %s\n", e.Count, e.Value)
	}
}
