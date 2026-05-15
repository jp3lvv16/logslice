// Package dropout implements probabilistic log line dropping.
// A Dropper randomly discards a configurable fraction of messages,
// useful for shedding load when log volume is too high.
package dropout

import (
	"fmt"
	"math/rand"
)

// Dropper randomly drops log lines based on a configured drop rate.
type Dropper struct {
	rate float64 // fraction of messages to drop, in [0.0, 1.0)
	rng  *rand.Rand
}

// New creates a Dropper that drops approximately rate*100% of messages.
// rate must be in the range [0.0, 1.0). A rate of 0.0 drops nothing;
// a rate of 0.9 drops ~90% of messages.
func New(rate float64, seed int64) (*Dropper, error) {
	if rate < 0.0 || rate >= 1.0 {
		return nil, fmt.Errorf("dropout: rate must be in [0.0, 1.0), got %v", rate)
	}
	return &Dropper{
		rate: rate,
		rng:  rand.New(rand.NewSource(seed)), //nolint:gosec
	}, nil
}

// Keep returns true if the message should be kept (i.e. not dropped).
// When rate is 0.0, Keep always returns true.
func (d *Dropper) Keep(_ []byte) bool {
	if d.rate == 0.0 {
		return true
	}
	return d.rng.Float64() >= d.rate
}

// Rate returns the configured drop rate.
func (d *Dropper) Rate() float64 {
	return d.rate
}
