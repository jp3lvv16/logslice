package main

import (
	"flag"
	"testing"
	"time"
)

func TestBuildThrottler_Disabled(t *testing.T) {
	th, err := buildThrottler(0, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th != nil {
		t.Fatal("expected nil throttler when max=0")
	}
}

func TestBuildThrottler_Valid(t *testing.T) {
	th, err := buildThrottler(50, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttler")
	}
}

func TestBuildThrottler_InvalidWindow(t *testing.T) {
	_, err := buildThrottler(10, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestRegisterThrottleFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var max int
	var window time.Duration
	registerThrottleFlags(fs, &max, &window)
	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if max != 0 {
		t.Errorf("expected default max=0, got %d", max)
	}
	if window != time.Second {
		t.Errorf("expected default window=1s, got %s", window)
	}
}

func TestRegisterThrottleFlags_CustomValues(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var max int
	var window time.Duration
	registerThrottleFlags(fs, &max, &window)
	if err := fs.Parse([]string{"--throttle-max=200", "--throttle-window=500ms"}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if max != 200 {
		t.Errorf("expected max=200, got %d", max)
	}
	if window != 500*time.Millisecond {
		t.Errorf("expected window=500ms, got %s", window)
	}
}
