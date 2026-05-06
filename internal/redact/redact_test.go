package redact

import (
	"encoding/json"
	"testing"
)

func TestNew_NoFields(t *testing.T) {
	_, err := New(nil, "")
	if err == nil {
		t.Fatal("expected error for empty field list")
	}
}

func TestNew_EmptyFieldName(t *testing.T) {
	_, err := New([]string{"ok", ""}, "")
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestNew_DefaultPlaceholder(t *testing.T) {
	r, err := New([]string{"password"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.placeholder != defaultPlaceholder {
		t.Errorf("got placeholder %q, want %q", r.placeholder, defaultPlaceholder)
	}
}

func TestApply_RedactsField(t *testing.T) {
	r, _ := New([]string{"token"}, "")
	input := json.RawMessage(`{"level":"info","token":"secret123","msg":"ok"}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["token"] != defaultPlaceholder {
		t.Errorf("token not redacted, got %q", obj["token"])
	}
	if obj["msg"] != "ok" {
		t.Errorf("msg field should be unchanged, got %q", obj["msg"])
	}
}

func TestApply_MultipleFields(t *testing.T) {
	r, _ := New([]string{"password", "ssn"}, "***")
	input := json.RawMessage(`{"user":"alice","password":"hunter2","ssn":"123-45-6789"}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["password"] != "***" {
		t.Errorf("password not redacted, got %q", obj["password"])
	}
	if obj["ssn"] != "***" {
		t.Errorf("ssn not redacted, got %q", obj["ssn"])
	}
	if obj["user"] != "alice" {
		t.Errorf("user should be unchanged, got %q", obj["user"])
	}
}

func TestApply_MissingField_NoError(t *testing.T) {
	r, _ := New([]string{"secret"}, "")
	input := json.RawMessage(`{"level":"debug","msg":"nothing sensitive"}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) == "" {
		t.Error("output should not be empty")
	}
}

func TestApply_NonJSONObject_PassThrough(t *testing.T) {
	r, _ := New([]string{"token"}, "")
	input := json.RawMessage(`not json at all`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected passthrough, got %q", string(out))
	}
}
