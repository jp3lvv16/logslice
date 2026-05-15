package normalize

import (
	"encoding/json"
	"testing"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_NoSpecs(t *testing.T) {
	n, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil Normalizer")
	}
}

func TestNew_ValidSpecs(t *testing.T) {
	_, err := New([]string{"level:lower", "msg:trim", "host:upper"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_MissingSeparator(t *testing.T) {
	_, err := New([]string{"levelonly"})
	if err == nil {
		t.Fatal("expected error for missing separator")
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New([]string{":lower"})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_UnknownMode(t *testing.T) {
	_, err := New([]string{"level:camel"})
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestApply_Lower(t *testing.T) {
	n, _ := New([]string{"level:lower"})
	out, err := n.Apply(raw(`{"level":"ERROR","msg":"boom"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if obj["level"] != "error" {
		t.Errorf("expected 'error', got %q", obj["level"])
	}
}

func TestApply_Upper(t *testing.T) {
	n, _ := New([]string{"level:upper"})
	out, _ := n.Apply(raw(`{"level":"warn"}`))
	var obj map[string]string
	json.Unmarshal(out, &obj)
	if obj["level"] != "WARN" {
		t.Errorf("expected 'WARN', got %q", obj["level"])
	}
}

func TestApply_Trim(t *testing.T) {
	n, _ := New([]string{"msg:trim"})
	out, _ := n.Apply(raw(`{"msg":"  hello world  "}`))
	var obj map[string]string
	json.Unmarshal(out, &obj)
	if obj["msg"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", obj["msg"])
	}
}

func TestApply_SkipsNonStringField(t *testing.T) {
	n, _ := New([]string{"count:lower"})
	input := raw(`{"count":42}`)
	out, err := n.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]json.RawMessage
	json.Unmarshal(out, &obj)
	if string(obj["count"]) != "42" {
		t.Errorf("expected numeric value preserved, got %s", obj["count"])
	}
}

func TestApply_SkipsMissingField(t *testing.T) {
	n, _ := New([]string{"level:lower"})
	input := raw(`{"msg":"hello"}`)
	out, err := n.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	json.Unmarshal(out, &obj)
	if _, ok := obj["level"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	n, _ := New(nil)
	input := raw(`{"level":"INFO"}`)
	out, _ := n.Apply(input)
	if string(out) != string(input) {
		t.Errorf("expected original message, got %s", out)
	}
}
