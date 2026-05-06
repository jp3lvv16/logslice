// Package schema validates that each log line conforms to a required set of
// field names, optionally enforcing their JSON types.
package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldType represents an expected JSON value type.
type FieldType string

const (
	TypeAny     FieldType = "any"
	TypeString  FieldType = "string"
	TypeNumber  FieldType = "number"
	TypeBoolean FieldType = "boolean"
)

// Rule describes a single required field and its expected type.
type Rule struct {
	Field string
	Type  FieldType
}

// Validator checks log lines against a set of field rules.
type Validator struct {
	rules []Rule
}

// New parses specs of the form "field" or "field:type" and returns a Validator.
// Valid types are: any, string, number, boolean.
func New(specs []string) (*Validator, error) {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		parts := strings.SplitN(s, ":", 2)
		rule := Rule{Field: parts[0], Type: TypeAny}
		if len(parts) == 2 {
			ft := FieldType(strings.ToLower(parts[1]))
			switch ft {
			case TypeAny, TypeString, TypeNumber, TypeBoolean:
				rule.Type = ft
			default:
				return nil, fmt.Errorf("schema: unknown type %q for field %q", parts[1], parts[0])
			}
		}
		rules = append(rules, rule)
	}
	return &Validator{rules: rules}, nil
}

// Validate returns an error if the raw JSON message is missing a required field
// or the field's value does not match the expected type.
func (v *Validator) Validate(raw json.RawMessage) error {
	if len(v.rules) == 0 {
		return nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return fmt.Errorf("schema: invalid JSON: %w", err)
	}
	for _, rule := range v.rules {
		val, ok := obj[rule.Field]
		if !ok {
			return fmt.Errorf("schema: missing required field %q", rule.Field)
		}
		if rule.Type == TypeAny {
			continue
		}
		if err := checkType(rule.Field, rule.Type, val); err != nil {
			return err
		}
	}
	return nil
}

func checkType(field string, ft FieldType, raw json.RawMessage) error {
	switch ft {
	case TypeString:
		var s string
		if json.Unmarshal(raw, &s) != nil {
			return fmt.Errorf("schema: field %q must be a string", field)
		}
	case TypeNumber:
		var n float64
		if json.Unmarshal(raw, &n) != nil {
			return fmt.Errorf("schema: field %q must be a number", field)
		}
	case TypeBoolean:
		var b bool
		if json.Unmarshal(raw, &b) != nil {
			return fmt.Errorf("schema: field %q must be a boolean", field)
		}
	}
	return nil
}
