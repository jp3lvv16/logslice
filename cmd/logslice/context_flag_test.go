package main

import (
	"encoding/json"
	"flag"
	"strings"
	"testing"
)

func TestBuildEnricher_Disabled(t *testing.T) {
	cf := contextFlags{Enrich: false}
	enrich := buildEnricher(cf)
	msg := json.RawMessage(`{"level":"info"}`)
	out, err := enrich(msg, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(out), "_line") {
		t.Error("_line should not be injected when enrich=false")
	}
}

func TestBuildEnricher_Enabled(t *testing.T) {
	cf := contextFlags{Enrich: true, Source: "test.log"}
	enrich := buildEnricher(cf)
	msg := json.RawMessage(`{"level":"info"}`)
	out, err := enrich(msg, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v, _ := obj["_line"].(float64); int(v) != 7 {
		t.Errorf("expected _line=7, got %v", obj["_line"])
	}
	if obj["_source"] != "test.log" {
		t.Errorf("expected _source=test.log, got %v", obj["_source"])
	}
}

func TestBuildStripper_Disabled(t *testing.T) {
	cf := contextFlags{Strip: false}
	strip := buildStripper(cf)
	msg := json.RawMessage(`{"_line":1,"level":"info"}`)
	out, err := strip(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "_line") {
		t.Error("_line should be preserved when strip-meta=false")
	}
}

func TestBuildStripper_Enabled(t *testing.T) {
	cf := contextFlags{Strip: true}
	strip := buildStripper(cf)
	msg := json.RawMessage(`{"_line":1,"_source":"f","_ingested_at":"x","level":"info"}`)
	out, err := strip(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, f := range []string{"_line", "_source", "_ingested_at"} {
		if strings.Contains(string(out), f) {
			t.Errorf("field %q should be stripped", f)
		}
	}
	if !strings.Contains(string(out), "level") {
		t.Error("level field should be preserved")
	}
}

func TestRegisterContextFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var cf contextFlags
	registerContextFlags(fs, &cf)
	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cf.Enrich {
		t.Error("enrich should default to false")
	}
	if cf.Strip {
		t.Error("strip-meta should default to false")
	}
	if cf.Source != "" {
		t.Error("source should default to empty")
	}
}
