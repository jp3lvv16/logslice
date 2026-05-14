// Package head limits processing to the first N log lines that pass all
// upstream filters. Once the limit is reached every subsequent line is
// dropped, allowing callers to implement a "top-N" view over a stream.
package head

import "errors"

// Limiter drops lines once a fixed count has been reached.
type Limiter struct {
	max  int
	seen int
}

// New creates a Limiter that keeps at most n lines.
// n == 0 means unlimited (Keep always returns true).
// A negative n is an error.
func New(n int) (*Limiter, error) {
	if n < 0 {
		return nil, errors.New("head: limit must be >= 0")
	}
	return &Limiter{max: n}, nil
}

// Keep returns true if the line should be forwarded downstream.
// It is safe to call Keep concurrently from a single goroutine; the
// Limiter is intentionally not goroutine-safe to avoid lock overhead
// in the single-threaded pipeline.
func (l *Limiter) Keep() bool {
	if l.max == 0 {
		return true
	}
	if l.seen >= l.max {
		return false
	}
	l.seen++
	return true
}

// Seen returns the number of lines that have been forwarded so far.
func (l *Limiter) Seen() int {
	return l.seen
}

// Done reports whether the limit has been reached and no further lines
// will ever be forwarded. Callers may use this to break out of a read
// loop early.
func (l *Limiter) Done() bool {
	if l.max == 0 {
		return false
	}
	return l.seen >= l.max
}
