// Package burst provides a sliding-window burst detector that flags when
// the number of log lines seen within a rolling duration exceeds a threshold.
package burst

import (
	"errors"
	"sync"
	"time"
)

// Detector tracks event timestamps in a sliding window and reports whether
// the current burst rate exceeds the configured threshold.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	timestamps []time.Time
	now       func() time.Time // injectable for testing
}

// New creates a Detector that fires when more than threshold events occur
// within the given window duration. Both window and threshold must be positive.
func New(window time.Duration, threshold int) (*Detector, error) {
	if window <= 0 {
		return nil, errors.New("burst: window must be positive")
	}
	if threshold <= 0 {
		return nil, errors.New("burst: threshold must be positive")
	}
	return &Detector{
		window:    window,
		threshold: threshold,
		now:       time.Now,
	}, nil
}

// Record registers one event and returns true if the number of events within
// the sliding window now exceeds the configured threshold.
func (d *Detector) Record() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	cutoff := now.Add(-d.window)

	// evict expired entries
	valid := d.timestamps[:0]
	for _, t := range d.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	d.timestamps = append(valid, now)

	return len(d.timestamps) > d.threshold
}

// Count returns the number of events currently inside the sliding window.
func (d *Detector) Count() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	cutoff := now.Add(-d.window)
	count := 0
	for _, t := range d.timestamps {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all recorded events.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.timestamps = d.timestamps[:0]
}
