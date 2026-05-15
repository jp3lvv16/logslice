// Package sample provides reservoir sampling for structured log lines.
// It collects up to N messages from an unbounded stream, giving each
// message an equal probability of appearing in the final output.
package sample

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

// Sampler holds a fixed-size reservoir of raw JSON messages.
type Sampler struct {
	reservoir []json.RawMessage
	size      int
	seen      int
	rng       *rand.Rand
}

// New creates a Sampler with a reservoir of the given size.
// size must be greater than zero.
func New(size int, seed int64) (*Sampler, error) {
	if size <= 0 {
		return nil, fmt.Errorf("sample: size must be > 0, got %d", size)
	}
	return &Sampler{
		reservoir: make([]json.RawMessage, 0, size),
		size:      size,
		rng:       rand.New(rand.NewSource(seed)),
	}, nil
}

// Add feeds a raw JSON message into the reservoir using Algorithm R.
// It returns true when the reservoir is full and a line was replaced.
func (s *Sampler) Add(msg json.RawMessage) {
	s.seen++
	if len(s.reservoir) < s.size {
		s.reservoir = append(s.reservoir, msg)
		return
	}
	j := s.rng.Intn(s.seen)
	if j < s.size {
		s.reservoir[j] = msg
	}
}

// Results returns the messages currently held in the reservoir.
// The order is not guaranteed.
func (s *Sampler) Results() []json.RawMessage {
	out := make([]json.RawMessage, len(s.reservoir))
	copy(out, s.reservoir)
	return out
}

// Seen returns the total number of messages that have been offered.
func (s *Sampler) Seen() int { return s.seen }

// Len returns the current number of messages in the reservoir.
func (s *Sampler) Len() int { return len(s.reservoir) }
