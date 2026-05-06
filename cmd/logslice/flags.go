package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type config struct {
	start      string
	end        string
	fields     []string
	filters    []string
	format     string
	sample     int
	tail       bool
	dedup      bool
	dedupWin   int
	rate       float64
	transforms []string
	truncLen   int
	truncField string
	flatten    bool
	aggregate  string
	redact     []string
	redactPH   string
	files      []string
}

func parseFlags() (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		start      = fs.String("start", "", "start of time range (RFC3339 or common layouts)")
		end        = fs.String("end", "", "end of time range")
		fields     = fs.String("fields", "", "comma-separated fields to include in output")
		filters    = fs.String("filter", "", "comma-separated field=pattern filter specs")
		format     = fs.String("format", "json", "output format: json|pretty|text")
		sample     = fs.Int("sample", 1, "emit every Nth matching line")
		tailF      = fs.Bool("tail", false, "follow file for new lines")
		dedup      = fs.Bool("dedup", false, "suppress duplicate log lines")
		dedupWin   = fs.Int("dedup-window", 1000, "number of recent lines to consider for dedup")
		rate       = fs.Float64("rate", 0, "max output lines per second (0 = unlimited)")
		transforms = fs.String("transform", "", "comma-separated transform specs (rename:old=new, drop:field, set:field=val)")
		truncLen   = fs.Int("trunc-len", 0, "max length for truncated field (0 = disabled)")
		truncField = fs.String("trunc-field", "msg", "field to apply truncation to")
		flatten    = fs.Bool("flatten", false, "flatten nested JSON objects")
		aggregate  = fs.String("aggregate", "", "field name to aggregate value counts for")
		redact     = fs.String("redact", "", "comma-separated field names whose values should be redacted")
		redactPH   = fs.String("redact-placeholder", "[REDACTED]", "replacement text used by --redact")
	)

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	cfg := &config{
		start:      *start,
		end:        *end,
		fields:     splitTrim(*fields),
		filters:    splitTrim(*filters),
		format:     *format,
		sample:     *sample,
		tail:       *tailF,
		dedup:      *dedup,
		dedupWin:   *dedupWin,
		rate:       *rate,
		transforms: splitTrim(*transforms),
		truncLen:   *truncLen,
		truncField: *truncField,
		flatten:    *flatten,
		aggregate:  *aggregate,
		redact:     splitTrim(*redact),
		redactPH:   *redactPH,
		files:      fs.Args(),
	}

	if cfg.format != "json" && cfg.format != "pretty" && cfg.format != "text" {
		return nil, fmt.Errorf("unknown format %q: must be json, pretty, or text", cfg.format)
	}
	return cfg, nil
}

func splitTrim(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
