// Package counter provides a simple line-count limiter that stops
// processing after a configurable maximum number of records have been emitted.
package counter

import "encoding/json"

// Counter tracks how many records have passed through and signals when the
// configured limit has been reached. A limit of zero means unlimited.
type Counter struct {
	max  int
	seen int
}

// New returns a Counter that allows at most max records through.
// Pass 0 (or any negative value) to disable limiting entirely.
func New(max int) *Counter {
	if max < 0 {
		max = 0
	}
	return &Counter{max: max}
}

// Allow returns true if the record should be forwarded and advances the
// internal counter. Once the limit is reached every subsequent call returns
// false. When max is 0 Allow always returns true.
func (c *Counter) Allow(msg json.RawMessage) bool {
	if c.max == 0 {
		return true
	}
	if c.seen >= c.max {
		return false
	}
	c.seen++
	return true
}

// Seen returns the number of records that have been allowed so far.
func (c *Counter) Seen() int {
	return c.seen
}

// Done reports whether the limit has been reached.
func (c *Counter) Done() bool {
	return c.max > 0 && c.seen >= c.max
}
