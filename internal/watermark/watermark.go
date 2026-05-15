// Package watermark tracks the high-water mark of a timestamp field across
// a stream of JSON log lines and emits a summary when the stream ends.
package watermark

import (
	"encoding/json"
	"time"
)

// Watermark tracks the earliest and latest timestamps seen in a log stream.
type Watermark struct {
	field string
	min   time.Time
	max   time.Time
	count int64
}

// New creates a Watermark that tracks timestamps stored in the given field.
// An empty field name returns an error.
func New(field string) (*Watermark, error) {
	if field == "" {
		return nil, errEmptyField
	}
	return &Watermark{field: field}, nil
}

// Record ingests a raw JSON log line and updates the min/max timestamps.
// Lines that are malformed or lack the target field are silently skipped.
func (w *Watermark) Record(line []byte) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(line, &m); err != nil {
		return
	}
	raw, ok := m[w.field]
	if !ok {
		return
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return
		}
	}
	w.count++
	if w.min.IsZero() || t.Before(w.min) {
		w.min = t
	}
	if w.max.IsZero() || t.After(w.max) {
		w.max = t
	}
}

// Min returns the earliest timestamp recorded, and whether any have been seen.
func (w *Watermark) Min() (time.Time, bool) {
	return w.min, !w.min.IsZero()
}

// Max returns the latest timestamp recorded, and whether any have been seen.
func (w *Watermark) Max() (time.Time, bool) {
	return w.max, !w.max.IsZero()
}

// Count returns the number of lines that contributed a valid timestamp.
func (w *Watermark) Count() int64 {
	return w.count
}

// Span returns the duration between the earliest and latest timestamps.
// Returns 0 if fewer than two timestamps have been recorded.
func (w *Watermark) Span() time.Duration {
	if w.min.IsZero() || w.max.IsZero() {
		return 0
	}
	return w.max.Sub(w.min)
}
