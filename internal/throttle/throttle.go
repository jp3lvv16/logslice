// Package throttle provides a token-bucket based line throttler that
// limits the number of log lines emitted per second across a pipeline.
package throttle

import (
	"encoding/json"
	"fmt"
	"time"
)

// Throttler drops log lines once the per-second budget is exhausted.
type Throttler struct {
	max     int
	window  time.Duration
	tokens  int
	reset   time.Time
	dropped int64
}

// New creates a Throttler that allows at most maxPerWindow lines inside
// each window duration. maxPerWindow <= 0 disables throttling.
func New(maxPerWindow int, window time.Duration) (*Throttler, error) {
	if window <= 0 {
		return nil, fmt.Errorf("throttle: window must be positive, got %s", window)
	}
	return &Throttler{
		max:    maxPerWindow,
		window: window,
		tokens: maxPerWindow,
		reset:  time.Now().Add(window),
	}, nil
}

// Allow returns true if the line represented by msg should be forwarded.
// It consumes one token from the current window budget.
func (t *Throttler) Allow(msg json.RawMessage) bool {
	if t.max <= 0 {
		return true
	}
	now := time.Now()
	if now.After(t.reset) {
		t.tokens = t.max
		t.reset = now.Add(t.window)
	}
	if t.tokens <= 0 {
		t.dropped++
		return false
	}
	t.tokens--
	return true
}

// Dropped returns the total number of lines dropped since creation.
func (t *Throttler) Dropped() int64 {
	return t.dropped
}
