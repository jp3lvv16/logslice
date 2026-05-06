package main

import (
	"encoding/json"
	"flag"

	ctx "github.com/user/logslice/internal/context"
)

// contextFlags holds the CLI flags that control metadata enrichment.
type contextFlags struct {
	Enrich bool
	Strip  bool
	Source string
}

// registerContextFlags attaches the context-related flags to fs.
func registerContextFlags(fs *flag.FlagSet, cf *contextFlags) {
	fs.BoolVar(&cf.Enrich, "enrich", false, "inject _line / _source / _ingested_at metadata into each record")
	fs.BoolVar(&cf.Strip, "strip-meta", false, "remove _line / _source / _ingested_at fields before output")
	fs.StringVar(&cf.Source, "source", "", "label written to _source when --enrich is active")
}

// buildEnricher returns a function that optionally enriches a JSON message
// with pipeline metadata, incrementing the line counter on each call.
func buildEnricher(cf contextFlags) func(json.RawMessage, int64) (json.RawMessage, error) {
	if !cf.Enrich {
		return func(msg json.RawMessage, _ int64) (json.RawMessage, error) { return msg, nil }
	}
	return func(msg json.RawMessage, line int64) (json.RawMessage, error) {
		meta := ctx.Meta{
			Line:   line,
			Source: cf.Source,
		}
		return ctx.Enrich(msg, meta)
	}
}

// buildStripper returns a function that optionally strips metadata fields.
func buildStripper(cf contextFlags) func(json.RawMessage) (json.RawMessage, error) {
	if !cf.Strip {
		return func(msg json.RawMessage) (json.RawMessage, error) { return msg, nil }
	}
	return ctx.Strip
}
