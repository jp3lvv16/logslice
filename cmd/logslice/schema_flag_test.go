package main

import (
	"testing"
)

func TestBuildValidator_NoSpecs(t *testing.T) {
	v, err := buildValidator(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != nil {
		t.Fatal("expected nil validator when no specs given")
	}
}

func TestBuildValidator_ValidSpecs(t *testing.T) {
	v, err := buildValidator([]string{"level:string", "ts:number"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil validator")
	}
}

func TestBuildValidator_InvalidSpec(t *testing.T) {
	_, err := buildValidator([]string{"level:unknowntype"})
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestSchemaFlag_Set(t *testing.T) {
	var f schemaFlag
	if err := f.Set("level:string"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := f.Set("ts:number"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(f))
	}
}

func TestSchemaFlag_String(t *testing.T) {
	f := schemaFlag{"level:string", "ts:number"}
	got := f.String()
	if got != "level:string,ts:number" {
		t.Fatalf("unexpected string: %q", got)
	}
}
