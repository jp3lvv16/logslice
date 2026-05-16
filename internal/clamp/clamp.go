// Package clamp truncates numeric field values to a configured [min, max] range.
// Values below min are raised to min; values above max are lowered to max.
// Fields that are absent or non-numeric are passed through unchanged.
package clamp

import (
	"encoding/json"
	"fmt"
	"math"
)

// Clamper rewrites a single numeric field so its value stays within [Min, Max].
type Clamper struct {
	field string
	min   float64
	max   float64
}

// New returns a Clamper that constrains field to [min, max].
// Returns an error when field is empty or min > max.
func New(field string, min, max float64) (*Clamper, error) {
	if field == "" {
		return nil, fmt.Errorf("clamp: field name must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("clamp: min (%g) must not exceed max (%g)", min, max)
	}
	return &Clamper{field: field, min: min, max: max}, nil
}

// Apply clamps the configured field in msg and returns the modified JSON.
// If the field is absent or its value is not a number the message is returned
// unchanged. Invalid JSON returns an error.
func (c *Clamper) Apply(msg json.RawMessage) (json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return nil, fmt.Errorf("clamp: unmarshal: %w", err)
	}

	raw, ok := obj[c.field]
	if !ok {
		return msg, nil
	}

	var v float64
	if err := json.Unmarshal(raw, &v); err != nil {
		// non-numeric — leave untouched
		return msg, nil
	}

	clamped := math.Max(c.min, math.Min(c.max, v))
	if clamped == v {
		return msg, nil
	}

	encoded, err := json.Marshal(clamped)
	if err != nil {
		return nil, fmt.Errorf("clamp: marshal value: %w", err)
	}
	obj[c.field] = encoded

	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("clamp: marshal object: %w", err)
	}
	return out, nil
}
