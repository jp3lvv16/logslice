package main

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/schema"
)

// schemaFlag is a repeatable CLI flag that accumulates field specs.
type schemaFlag []string

func (f *schemaFlag) String() string { return strings.Join(*f, ",") }
func (f *schemaFlag) Set(v string) error {
	*f = append(*f, strings.TrimSpace(v))
	return nil
}

// buildValidator constructs a schema.Validator from the collected specs,
// or returns nil if no specs were provided.
func buildValidator(specs []string) (*schema.Validator, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	v, err := schema.New(specs)
	if err != nil {
		return nil, fmt.Errorf("--schema: %w", err)
	}
	return v, nil
}
