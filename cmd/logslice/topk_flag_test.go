package main

import (
	"flag"
	"testing"
)

func TestBuildTopKTracker_Disabled(t *testing.T) {
	tr, err := buildTopKTracker("", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr != nil {
		t.Error("expected nil tracker when field is empty")
	}
}

func TestBuildTopKTracker_Valid(t *testing.T) {
	tr, err := buildTopKTracker("level", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestBuildTopKTracker_InvalidK(t *testing.T) {
	_, err := buildTopKTracker("level", 0)
	if err == nil {
		t.Fatal("expected error for k=0")
	}
}

func TestRegisterTopKFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var field string
	var k int
	registerTopKFlags(fs, &field, &k)

	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if field != "" {
		t.Errorf("expected empty default field, got %q", field)
	}
	if k != 10 {
		t.Errorf("expected default k=10, got %d", k)
	}
}

func TestRegisterTopKFlags_CustomValues(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var field string
	var k int
	registerTopKFlags(fs, &field, &k)

	if err := fs.Parse([]string{"--topk-field", "service", "--topk-k", "3"}); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if field != "service" {
		t.Errorf("expected field=service, got %q", field)
	}
	if k != 3 {
		t.Errorf("expected k=3, got %d", k)
	}
}
