package main

import (
	"flag"
	"fmt"
	"os"
)

// config holds all parsed CLI flags for a logslice run.
type config struct {
	start   string
	end     string
	fields  string
	filters []string
	format  string
	sample  int
	files   []string
}

// parseFlags parses os.Args and returns a populated config.
func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		start   = fs.String("start", "", "start of time range (RFC3339 or common log formats)")
		end     = fs.String("end", "", "end of time range (RFC3339 or common log formats)")
		fields  = fs.String("fields", "", "comma-separated list of fields to include in output")
		format  = fs.String("format", "json", "output format: json, pretty, or text")
		sample  = fs.Int("sample", 1, "emit every Nth matching line (>=1)")
		filter  multiFlag
	)

	fs.Var(&filter, "filter", "field=pattern filter (repeatable); e.g. -filter level=error")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *sample < 1 {
		return nil, fmt.Errorf("--sample must be >= 1")
	}

	return &config{
		start:   *start,
		end:     *end,
		fields:  *fields,
		filters: []string(filter),
		format:  *format,
		sample:  *sample,
		files:   fs.Args(),
	}, nil
}

// multiFlag is a flag.Value that accumulates repeated flag values.
type multiFlag []string

func (m *multiFlag) String() string {
	if m == nil {
		return ""
	}
	return fmt.Sprintf("%v", []string(*m))
}

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}
