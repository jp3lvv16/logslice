// Package window provides a sliding time-window buffer that groups
// consecutive log lines falling within a configurable duration.
package window

import (
	"encoding/json"
	"time"
)

// Group holds a slice of raw JSON messages that share the same time window.
type Group struct {
	Start    time.Time
	Messages []json.RawMessage
}

// Window accumulates log lines into fixed-duration buckets.
type Window struct {
	size    time.Duration
	current *Group
}

// New creates a Window with the given bucket size.
// A size of zero means every message forms its own group.
func New(size time.Duration) (*Window, error) {
	if size < 0 {
		return nil, errNegativeSize
	}
	return &Window{size: size}, nil
}

// Add places msg into the appropriate bucket based on ts.
// If the message belongs to a new bucket the completed previous group is
// returned together with a true flush flag; otherwise flush is false.
func (w *Window) Add(ts time.Time, msg json.RawMessage) (flushed *Group, flush bool) {
	if w.size == 0 {
		g := &Group{Start: ts, Messages: []json.RawMessage{msg}}
		return g, true
	}

	if w.current == nil {
		w.current = &Group{Start: ts}
	}

	deadline := w.current.Start.Add(w.size)
	if !ts.Before(deadline) {
		flushed = w.current
		w.current = &Group{Start: ts}
		w.current.Messages = append(w.current.Messages, msg)
		return flushed, true
	}

	w.current.Messages = append(w.current.Messages, msg)
	return nil, false
}

// Flush returns any buffered group and resets internal state.
// Call after the input stream is exhausted to drain the last bucket.
func (w *Window) Flush() *Group {
	g := w.current
	w.current = nil
	return g
}
