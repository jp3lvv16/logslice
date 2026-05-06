// Package sampler provides rate-based log line sampling for logslice.
// It allows users to process only every Nth matching log line, which is
// useful for reducing output volume when working with high-frequency logs.
package sampler

import "fmt"

// Sampler tracks line counts and decides whether a given line should be
// emitted based on a configured sample rate.
type Sampler struct {
	rate    int
	counter int
}

// New creates a Sampler that emits every nth line.
// A rate of 1 (or less) means every line is emitted (no sampling).
func New(rate int) (*Sampler, error) {
	if rate < 1 {
		return nil, fmt.Errorf("sampler: rate must be >= 1, got %d", rate)
	}
	return &Sampler{rate: rate}, nil
}

// Keep increments the internal counter and returns true if the current
// line should be kept according to the configured sample rate.
func (s *Sampler) Keep() bool {
	s.counter++
	return s.counter%s.rate == 0
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() int {
	return s.rate
}

// Reset resets the internal counter to zero.
func (s *Sampler) Reset() {
	s.counter = 0
}

// Count returns the total number of lines seen so far.
func (s *Sampler) Count() int {
	return s.counter
}
