// Package reader provides utilities for reading structured log lines
// from an io.Reader, parsing each line as JSON, and extracting timestamps.
package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Line represents a single parsed log line.
type Line struct {
	Raw       string
	Fields    map[string]interface{}
	Timestamp time.Time
}

// TimestampKeys is the ordered list of JSON field names tried when
// extracting a timestamp from a log line.
var TimestampKeys = []string{"time", "timestamp", "ts", "@timestamp"}

// Read scans r line-by-line, parses each line as JSON, and sends
// successfully parsed lines to the returned channel. Malformed JSON
// lines are skipped. The channel is closed when r is exhausted or an
// unrecoverable read error occurs.
func Read(r io.Reader) <-chan Line {
	ch := make(chan Line)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			raw := scanner.Text()
			if raw == "" {
				continue
			}
			fields := make(map[string]interface{})
			if err := json.Unmarshal([]byte(raw), &fields); err != nil {
				continue
			}
			ts := extractTimestamp(fields)
			ch <- Line{Raw: raw, Fields: fields, Timestamp: ts}
		}
	}()
	return ch
}

// extractTimestamp attempts to parse a timestamp from well-known fields.
func extractTimestamp(fields map[string]interface{}) time.Time {
	for _, key := range TimestampKeys {
		v, ok := fields[key]
		if !ok {
			continue
		}
		switch val := v.(type) {
		case string:
			if t, err := parseFlexible(val); err == nil {
				return t
			}
		case float64:
			// Unix seconds (possibly fractional)
			return time.Unix(int64(val), 0).UTC()
		}
	}
	return time.Time{}
}

var timeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
}

func parseFlexible(s string) (time.Time, error) {
	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("reader: cannot parse time %q", s)
}
