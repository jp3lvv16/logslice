package main

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/mask"
)

// maskFlag is a repeatable CLI flag that accumulates mask specs.
type maskFlag []string

func (f *maskFlag) String() string {
	return strings.Join(*f, ", ")
}

func (f *maskFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
}

// buildMasker constructs a *mask.Masker from the collected specs.
// Returns nil and no error when specs is empty.
func buildMasker(specs []string) (*mask.Masker, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	m, err := mask.New(specs)
	if err != nil {
		return nil, fmt.Errorf("--mask: %w", err)
	}
	return m, nil
}
