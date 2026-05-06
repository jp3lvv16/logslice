package main

import (
	"bytes"
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.format != "json" {
		t.Errorf("expected default format json, got %q", cfg.format)
	}
	if cfg.file != "" {
		t.Errorf("expected empty file, got %q", cfg.file)
	}
}

func TestParseFlags_TimeRange(t *testing.T) {
	cfg, err := parseFlags([]string{"-from", "2024-01-01T00:00:00Z", "-to", "2024-01-02T00:00:00Z"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.from != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected from: %q", cfg.from)
	}
	if cfg.to != "2024-01-02T00:00:00Z" {
		t.Errorf("unexpected to: %q", cfg.to)
	}
}

func TestParseFlags_Fields(t *testing.T) {
	cfg, err := parseFlags([]string{"-fields", "level, msg, ts"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(cfg.fields))
	}
	if cfg.fields[1] != "msg" {
		t.Errorf("expected msg, got %q", cfg.fields[1])
	}
}

func TestParseFlags_Filters(t *testing.T) {
	cfg, err := parseFlags([]string{"-filter", "level=error,service=api"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(cfg.filters))
	}
}

func TestParseFlags_Format(t *testing.T) {
	for _, fmt := range []string{"json", "pretty", "text"} {
		cfg, err := parseFlags([]string{"-format", fmt}, &bytes.Buffer{})
		if err != nil {
			t.Fatalf("unexpected error for format %q: %v", fmt, err)
		}
		if cfg.format != fmt {
			t.Errorf("expected format %q, got %q", fmt, cfg.format)
		}
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := parseFlags([]string{"-unknown", "value"}, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}
