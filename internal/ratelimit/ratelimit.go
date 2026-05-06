// Package ratelimit provides a token-bucket style rate limiter for log line
// throughput, allowing callers to cap how many lines per second are emitted.
package ratelimit

import (
	"errors"
	"time"
)

// Limiter controls how many log lines may pass through per second.
type Limiter struct {
	tokens   float64
	max      float64
	rate     float64 // tokens per nanosecond
	lastTick time.Time
	now      func() time.Time
}

// New returns a Limiter that allows at most linesPerSec lines per second.
// linesPerSec must be greater than zero.
func New(linesPerSec float64) (*Limiter, error) {
	if linesPerSec <= 0 {
		return nil, errors.New("ratelimit: linesPerSec must be > 0")
	}
	return &Limiter{
		tokens:   linesPerSec,
		max:      linesPerSec,
		rate:     linesPerSec / float64(time.Second),
		lastTick: time.Now(),
		now:      time.Now,
	}, nil
}

// Allow reports whether the next log line should be allowed through.
// It refills tokens based on elapsed time since the last call.
func (l *Limiter) Allow() bool {
	now := l.now()
	elapsed := now.Sub(l.lastTick)
	l.lastTick = now

	l.tokens += float64(elapsed) * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// Unlimited returns a Limiter that always allows every line through.
func Unlimited() *Limiter {
	return &Limiter{
		tokens:   1,
		max:      1,
		rate:     0,
		lastTick: time.Now(),
		now:      time.Now,
	}
}

// Allow on an unlimited limiter always returns true.
func (l *Limiter) isUnlimited() bool {
	return l.rate == 0
}
