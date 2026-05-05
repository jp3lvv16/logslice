package timerange

import (
	"fmt"
	"time"
)

// TimeRange represents an inclusive start/end time window for filtering log entries.
type TimeRange struct {
	Start *time.Time
	End   *time.Time
}

// DefaultLayout is the default time format used for parsing time range arguments.
const DefaultLayout = time.RFC3339

// Parse parses start and end time strings into a TimeRange.
// Either start or end may be empty, resulting in an open-ended range.
func Parse(start, end string) (TimeRange, error) {
	var tr TimeRange

	if start != "" {
		t, err := parseTime(start)
		if err != nil {
			return tr, fmt.Errorf("invalid start time %q: %w", start, err)
		}
		tr.Start = &t
	}

	if end != "" {
		t, err := parseTime(end)
		if err != nil {
			return tr, fmt.Errorf("invalid end time %q: %w", end, err)
		}
		tr.End = &t
	}

	if tr.Start != nil && tr.End != nil && tr.End.Before(*tr.Start) {
		return tr, fmt.Errorf("end time must not be before start time")
	}

	return tr, nil
}

// Contains reports whether t falls within the TimeRange (inclusive).
func (tr TimeRange) Contains(t time.Time) bool {
	if tr.Start != nil && t.Before(*tr.Start) {
		return false
	}
	if tr.End != nil && t.After(*tr.End) {
		return false
	}
	return true
}

// IsZero reports whether the TimeRange has no bounds set.
func (tr TimeRange) IsZero() bool {
	return tr.Start == nil && tr.End == nil
}

var layouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

func parseTime(s string) (time.Time, error) {
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised time format (try RFC3339, e.g. 2006-01-02T15:04:05Z)")
}
