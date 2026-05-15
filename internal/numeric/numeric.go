// Package numeric provides field-based numeric range filtering for structured
// log lines. A Filter keeps only messages where a named numeric field falls
// within an inclusive [Min, Max] interval.
package numeric

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

// Filter rejects log lines whose numeric field value lies outside the
// configured range.
type Filter struct {
	field string
	min   float64
	max   float64
}

// New creates a Filter that keeps messages where msg[field] ∈ [min, max].
// Pass math.Inf(-1) or math.Inf(1) for an open-ended bound.
func New(field string, min, max float64) (*Filter, error) {
	if field == "" {
		return nil, errors.New("numeric: field name must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("numeric: min (%g) must not exceed max (%g)", min, max)
	}
	return &Filter{field: field, min: min, max: max}, nil
}

// Keep returns true when the raw JSON message contains the target field and
// its numeric value is within [min, max]. Lines that are malformed, missing
// the field, or carry a non-numeric value are dropped (returns false).
func (f *Filter) Keep(line []byte) bool {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(line, &obj); err != nil {
		return false
	}
	raw, ok := obj[f.field]
	if !ok {
		return false
	}
	var v float64
	if err := json.Unmarshal(raw, &v); err != nil {
		return false
	}
	return v >= f.min && v <= f.max
}

// Bounds returns the configured [min, max] interval.
func (f *Filter) Bounds() (min, max float64) {
	return f.min, f.max
}

// Field returns the name of the field being inspected.
func (f *Filter) Field() string {
	return f.field
}

// String returns a human-readable description of the filter.
func (f *Filter) String() string {
	lo, hi := formatBound(f.min, false), formatBound(f.max, true)
	return fmt.Sprintf("numeric: %s in [%s, %s]", f.field, lo, hi)
}

func formatBound(v float64, upper bool) string {
	switch {
	case upper && math.IsInf(v, 1):
		return "+Inf"
	case !upper && math.IsInf(v, -1):
		return "-Inf"
	default:
		return fmt.Sprintf("%g", v)
	}
}
