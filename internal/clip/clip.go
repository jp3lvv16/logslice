// Package clip provides a sliding-window byte-budget limiter that drops
// log messages once the cumulative encoded size of emitted lines within a
// rolling time window exceeds a configured maximum.
package clip

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// Clipper tracks byte usage inside a sliding window and decides whether a
// message should be forwarded or silently dropped.
type Clipper struct {
	mu      sync.Mutex
	maxBytes int64
	window   time.Duration
	buckets  []bucket
}

type bucket struct {
	at    time.Time
	bytes int64
}

// New creates a Clipper that allows at most maxBytes total encoded bytes to
// pass through within the given rolling window duration.
// maxBytes <= 0 disables clipping (all messages are allowed).
func New(maxBytes int64, window time.Duration) (*Clipper, error) {
	if window <= 0 {
		return nil, errors.New("clip: window must be positive")
	}
	return &Clipper{
		maxBytes: maxBytes,
		window:   window,
	}, nil
}

// Allow reports whether msg should be forwarded given the current byte budget.
// It advances the internal clock and evicts stale buckets on every call.
func (c *Clipper) Allow(msg json.RawMessage) bool {
	if c.maxBytes <= 0 {
		return true
	}

	now := time.Now()
	size := int64(len(msg))

	c.mu.Lock()
	defer c.mu.Unlock()

	c.evict(now)

	var used int64
	for _, b := range c.buckets {
		used += b.bytes
	}

	if used+size > c.maxBytes {
		return false
	}

	c.buckets = append(c.buckets, bucket{at: now, bytes: size})
	return true
}

// Used returns the number of bytes consumed in the current window.
func (c *Clipper) Used() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict(time.Now())
	var total int64
	for _, b := range c.buckets {
		total += b.bytes
	}
	return total
}

func (c *Clipper) evict(now time.Time) {
	cutoff := now.Add(-c.window)
	i := 0
	for i < len(c.buckets) && c.buckets[i].at.Before(cutoff) {
		i++
	}
	c.buckets = c.buckets[i:]
}
