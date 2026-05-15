package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/user/logslice/internal/fieldmap"
)

// fieldmapFlag is a repeatable flag that accumulates "old:new" rename specs.
type fieldmapFlag []string

func (f *fieldmapFlag) String() string { return strings.Join(*f, ",") }
func (f *fieldmapFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
}

// registerFieldmapFlags adds the -rename flag to fs and returns a pointer to
// the accumulated specs slice.
func registerFieldmapFlags(fs *flag.FlagSet) *fieldmapFlag {
	var specs fieldmapFlag
	fs.Var(&specs, "rename", `rename a JSON field: "old:new" (repeatable)`)
	return &specs
}

// buildMapper constructs a *fieldmap.Mapper from the collected specs.
// Returns nil, nil when no specs were provided so callers can skip the step.
func buildMapper(specs []string) (*fieldmap.Mapper, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	m, err := fieldmap.New(specs)
	if err != nil {
		return nil, fmt.Errorf("--rename: %w", err)
	}
	return m, nil
}
