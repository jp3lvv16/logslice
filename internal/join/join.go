// Package join correlates log lines from multiple sources by a shared key.
// Lines arriving within the join window that share the same key value are
// merged into a single JSON object before being emitted.
package join

import (
	"encoding/json"
	"fmt"
	"time"
)

// Joiner buffers log lines and merges those that share a key within a window.
type Joiner struct {
	key    string
	window time.Duration
	buf    map[string]*entry
}

type entry struct {
	merged map[string]json.RawMessage
	expiry time.Time
}

// New creates a Joiner that correlates records by key within the given window.
// key must be non-empty and window must be positive.
func New(key string, window time.Duration) (*Joiner, error) {
	if key == "" {
		return nil, fmt.Errorf("join: key must not be empty")
	}
	if window <= 0 {
		return nil, fmt.Errorf("join: window must be positive")
	}
	return &Joiner{
		key:    key,
		window: window,
		buf:    make(map[string]*entry),
	}, nil
}

// Add ingests a raw JSON log line. It returns a merged JSON object if the
// entry has been seen before within the window, otherwise nil (buffered).
// Expired entries are evicted before processing.
func (j *Joiner) Add(line []byte) ([]byte, error) {
	now := time.Now()
	j.evict(now)

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(line, &fields); err != nil {
		return nil, fmt.Errorf("join: unmarshal: %w", err)
	}

	kv, ok := fields[j.key]
	if !ok {
		return nil, nil
	}
	var keyStr string
	if err := json.Unmarshal(kv, &keyStr); err != nil {
		return nil, fmt.Errorf("join: key value must be a string: %w", err)
	}

	if e, exists := j.buf[keyStr]; exists {
		for k, v := range fields {
			e.merged[k] = v
		}
		delete(j.buf, keyStr)
		out, err := json.Marshal(e.merged)
		if err != nil {
			return nil, fmt.Errorf("join: marshal: %w", err)
		}
		return out, nil
	}

	j.buf[keyStr] = &entry{
		merged: fields,
		expiry: now.Add(j.window),
	}
	return nil, nil
}

// Flush returns all buffered (unmatched) entries and clears the buffer.
func (j *Joiner) Flush() ([][]byte, error) {
	var out [][]byte
	for _, e := range j.buf {
		b, err := json.Marshal(e.merged)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	j.buf = make(map[string]*entry)
	return out, nil
}

func (j *Joiner) evict(now time.Time) {
	for k, e := range j.buf {
		if now.After(e.expiry) {
			delete(j.buf, k)
		}
	}
}
