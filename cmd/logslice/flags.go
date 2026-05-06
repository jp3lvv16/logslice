package main

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type config struct {
	from    string
	to      string
	format  string
	file    string
	fields  []string
	filters []string
}

func parseFlags(args []string, stderr io.Writer) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		from    = fs.String("from", "", "start of time range (RFC3339 or common layouts)")
		to      = fs.String("to", "", "end of time range (RFC3339 or common layouts)")
		format  = fs.String("format", "json", "output format: json, pretty, text")
		file    = fs.String("file", "", "input log file (defaults to stdin)")
		fields  = fs.String("fields", "", "comma-separated list of fields to include in output")
		filterF = fs.String("filter", "", "comma-separated field=pattern filter specs")
	)

	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: logslice [options]")
		fmt.Fprintln(stderr, "\nOptions:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &config{
		from:   *from,
		to:     *to,
		format: *format,
		file:   *file,
	}

	if *fields != "" {
		for _, f := range strings.Split(*fields, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				cfg.fields = append(cfg.fields, f)
			}
		}
	}

	if *filterF != "" {
		for _, f := range strings.Split(*filterF, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				cfg.filters = append(cfg.filters, f)
			}
		}
	}

	return cfg, nil
}
