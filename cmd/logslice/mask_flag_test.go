package main

import (
	"testing"
)

func TestBuildMasker_NoSpecs(t *testing.T) {
	m, err := buildMasker(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m != nil {
		t.Error("expected nil masker for empty specs")
	}
}

func TestBuildMasker_ValidSpecs(t *testing.T) {
	m, err := buildMasker([]string{"token/[A-Z]+/***"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil masker")
	}
}

func TestBuildMasker_InvalidSpec(t *testing.T) {
	_, err := buildMasker([]string{"bad-spec"})
	if err == nil {
		t.Fatal("expected error for invalid spec")
	}
}

func TestBuildMasker_InvalidRegex(t *testing.T) {
	_, err := buildMasker([]string{"field/(/mask"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestMaskFlag_Set(t *testing.T) {
	var f maskFlag
	if err := f.Set("email/[^@]+@/***@"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := f.Set("token/.*/REDACTED"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if len(f) != 2 {
		t.Errorf("len = %d, want 2", len(f))
	}
}

func TestMaskFlag_String(t *testing.T) {
	f := maskFlag{"a/b/c", "d/e/f"}
	s := f.String()
	if s != "a/b/c, d/e/f" {
		t.Errorf("String() = %q", s)
	}
}
