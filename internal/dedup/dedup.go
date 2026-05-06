// Package dedup provides line-level deduplication for structured log streams.
// It tracks a sliding window of message fingerprints and suppresses repeated
// entries within that window, reducing noise in high-volume log output.
package dedup

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// Deduplicator tracks seen log line fingerprints within a fixed-size window.
type Deduplicator struct {
	window int
	seen   []string
	index  int
	set    map[string]struct{}
}

// New creates a Deduplicator with the given window size.
// window is the number of recent fingerprints to remember.
// A window of 0 disables deduplication (all lines are kept).
func New(window int) *Deduplicator {
	if window <= 0 {
		return &Deduplicator{window: 0}
	}
	return &Deduplicator{
		window: window,
		seen:   make([]string, window),
		set:    make(map[string]struct{}, window),
	}
}

// IsDuplicate returns true if the raw JSON line has been seen within the
// current window. It always adds the line to the window on first sight.
func (d *Deduplicator) IsDuplicate(line []byte) bool {
	if d.window == 0 {
		return false
	}

	fp := fingerprint(line)

	if _, exists := d.set[fp]; exists {
		return true
	}

	// Evict oldest entry if window is full.
	old := d.seen[d.index]
	if old != "" {
		delete(d.set, old)
	}

	d.seen[d.index] = fp
	d.set[fp] = struct{}{}
	d.index = (d.index + 1) % d.window

	return false
}

// fingerprint returns a stable hash for a JSON object, ignoring key order.
func fingerprint(line []byte) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(line, &obj); err != nil {
		// Fall back to raw byte hash for non-JSON lines.
		h := sha256.Sum256(line)
		return fmt.Sprintf("%x", h)
	}

	// Re-marshal with sorted keys (encoding/json sorts map keys).
	norm, err := json.Marshal(obj)
	if err != nil {
		h := sha256.Sum256(line)
		return fmt.Sprintf("%x", h)
	}

	h := sha256.Sum256(norm)
	return fmt.Sprintf("%x", h)
}
