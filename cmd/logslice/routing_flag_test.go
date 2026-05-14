package main

import (
	"flag"
	"testing"
)

func TestBuildRouter_NoSpecs(t *testing.T) {
	r, err := buildRouter(nil, "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != nil {
		t.Error("expected nil router when no specs given")
	}
}

func TestBuildRouter_ValidSpecs(t *testing.T) {
	r, err := buildRouter([]string{"level=error:errors", "level=warn:warnings"}, "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil router")
	}
	dest := r.Route([]byte(`{"level":"error"}`))
	if dest != "errors" {
		t.Errorf("expected 'errors', got %q", dest)
	}
}

func TestBuildRouter_InvalidSpec(t *testing.T) {
	_, err := buildRouter([]string{"bad-spec"}, "default")
	if err == nil {
		t.Fatal("expected error for invalid spec")
	}
}

func TestRoutingFlag_Set(t *testing.T) {
	var specs routingFlag
	if err := specs.Set("level=error:errors"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := specs.Set("level=warn:warnings"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 2 {
		t.Errorf("expected 2 specs, got %d", len(specs))
	}
}

func TestRoutingFlag_String(t *testing.T) {
	specs := routingFlag{"level=error:errors", "level=warn:warnings"}
	got := specs.String()
	if got != "level=error:errors,level=warn:warnings" {
		t.Errorf("unexpected String(): %q", got)
	}
}

func TestRegisterRoutingFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var specs routingFlag
	var defaultDest string
	registerRoutingFlags(fs, &specs, &defaultDest)
	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if defaultDest != "default" {
		t.Errorf("expected default dest 'default', got %q", defaultDest)
	}
}
