package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// config holds all parsed CLI flags.
type config struct {
	files      []string
	start      string
	end        string
	fields     []string
	filters    []string
	format     string
	sample     int
	dedupKeys  []string
	dedupWin   int
	rateLimit  float64
	highlight  bool
	tailFollow bool
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		start     = fs.String("start", "", "start of time range (RFC3339 or common layouts)")
		end       = fs.String("end", "", "end of time range")
		fields    = fs.String("fields", "", "comma-separated list of fields to include in output")
		filters   = fs.String("filter", "", "comma-separated field=pattern filter specs")
		format    = fs.String("format", "json", "output format: json, pretty, text")
		sample    = fs.Int("sample", 1, "keep every Nth line (1 = keep all)")
		dedupKeys = fs.String("dedup", "", "comma-separated fields to use as dedup key")
		dedupWin  = fs.Int("dedup-window", 1000, "number of recent fingerprints to track for dedup")
		rateLimit = fs.Float64("rate", 0, "max lines/sec to output (0 = unlimited)")
		highlight = fs.Bool("highlight", false, "colorize level field in terminal output")
		follow    = fs.Bool("follow", false, "follow file for new lines (like tail -f)")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *format != "json" && *format != "pretty" && *format != "text" {
		return nil, fmt.Errorf("unknown format %q: must be json, pretty or text", *format)
	}

	cfg := &config{
		files:      fs.Args(),
		start:      *start,
		end:        *end,
		format:     *format,
		sample:     *sample,
		dedupWin:   *dedupWin,
		rateLimit:  *rateLimit,
		highlight:  *highlight,
		tailFollow: *follow,
	}

	if *fields != "" {
		cfg.fields = splitTrim(*fields)
	}
	if *filters != "" {
		cfg.filters = splitTrim(*filters)
	}
	if *dedupKeys != "" {
		cfg.dedupKeys = splitTrim(*dedupKeys)
	}

	return cfg, nil
}

func splitTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
