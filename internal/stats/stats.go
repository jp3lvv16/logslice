// Package stats tracks processing metrics for a logslice run.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Counter accumulates counts of lines seen during processing.
type Counter struct {
	Total    int
	Matched  int
	Skipped  int
	Malformed int
	start    time.Time
}

// New returns a new Counter with the start time set to now.
func New() *Counter {
	return &Counter{start: time.Now()}
}

// IncTotal increments the total lines seen.
func (c *Counter) IncTotal() { c.Total++ }

// IncMatched increments the matched (output) line count.
func (c *Counter) IncMatched() { c.Matched++ }

// IncSkipped increments the skipped (out-of-range / filtered) line count.
func (c *Counter) IncSkipped() { c.Skipped++ }

// IncMalformed increments the malformed (unparseable) line count.
func (c *Counter) IncMalformed() { c.Malformed++ }

// Elapsed returns the duration since the counter was created.
func (c *Counter) Elapsed() time.Duration {
	return time.Since(c.start)
}

// Print writes a human-readable summary to w.
func (c *Counter) Print(w io.Writer) {
	fmt.Fprintf(w, "--- stats ---\n")
	fmt.Fprintf(w, "total:     %d\n", c.Total)
	fmt.Fprintf(w, "matched:   %d\n", c.Matched)
	fmt.Fprintf(w, "skipped:   %d\n", c.Skipped)
	fmt.Fprintf(w, "malformed: %d\n", c.Malformed)
	fmt.Fprintf(w, "elapsed:   %s\n", c.Elapsed().Round(time.Millisecond))
}
