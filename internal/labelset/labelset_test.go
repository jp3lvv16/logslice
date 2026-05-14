package labelset

import (
	"encoding/json"
	"testing"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_ValidSpecs(t *testing.T) {
	l, err := New([]string{"env=prod", "region=us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(l.labels))
	}
}

func TestNew_EmptySpecs(t *testing.T) {
	l, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.labels) != 0 {
		t.Fatal("expected empty label map")
	}
}

func TestNew_MissingSeparator(t *testing.T) {
	_, err := New([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestApply_AddsLabels(t *testing.T) {
	l, _ := New([]string{"env=prod", "region=eu"})
	out, err := l.Apply(raw(`{"msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", obj["env"])
	}
	if obj["region"] != "eu" {
		t.Errorf("expected region=eu, got %q", obj["region"])
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	l, _ := New([]string{"env=prod"})
	out, err := l.Apply(raw(`{"env":"dev"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	json.Unmarshal(out, &obj)
	if obj["env"] != "dev" {
		t.Errorf("existing value should be preserved, got %q", obj["env"])
	}
}

func TestApply_NoLabels_ReturnsUnchanged(t *testing.T) {
	l, _ := New(nil)
	input := raw(`{"msg":"hi"}`)
	out, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output")
	}
}

func TestApply_InvalidJSON_ReturnsOriginal(t *testing.T) {
	l, _ := New([]string{"env=prod"})
	input := raw(`not-json`)
	out, err := l.Apply(input)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if string(out) != string(input) {
		t.Errorf("expected original bytes on error")
	}
}
