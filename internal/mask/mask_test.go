package mask

import (
	"encoding/json"
	"testing"
)

func raw(t *testing.T, m map[string]string) []byte {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return b
}

func TestNew_ValidSpec(t *testing.T) {
	_, err := New([]string{"email/[^@]+@/***@"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidSpec(t *testing.T) {
	_, err := New([]string{"bad-spec"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := New([]string{"field/(/mask"})
	if err == nil {
		t.Fatal("expected error for bad regex")
	}
}

func TestApply_MasksMatchingField(t *testing.T) {
	m, _ := New([]string{"token/[A-Za-z0-9]{8,}/***REDACTED***"})
	input := raw(t, map[string]string{"token": "abc12345xyz", "msg": "hello"})
	out := m.Apply(input)

	var obj map[string]string
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if obj["token"] != "***REDACTED***" {
		t.Errorf("token = %q, want ***REDACTED***", obj["token"])
	}
	if obj["msg"] != "hello" {
		t.Errorf("msg = %q, want hello", obj["msg"])
	}
}

func TestApply_NoMatch_Unchanged(t *testing.T) {
	m, _ := New([]string{"token/[0-9]{10,}/***"})
	input := raw(t, map[string]string{"token": "short"})
	out := m.Apply(input)
	if string(out) != string(input) {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	m, _ := New([]string{"secret/.*/**"})
	input := raw(t, map[string]string{"msg": "no secret here"})
	out := m.Apply(input)
	if string(out) != string(input) {
		t.Errorf("expected unchanged, got %s", out)
	}
}

func TestApply_NonJSON_Passthrough(t *testing.T) {
	m, _ := New([]string{"field/.*/mask"})
	input := []byte("not json at all")
	out := m.Apply(input)
	if string(out) != string(input) {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_NoRules_Passthrough(t *testing.T) {
	m, _ := New(nil)
	input := raw(t, map[string]string{"k": "v"})
	out := m.Apply(input)
	if string(out) != string(input) {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_MultipleRules(t *testing.T) {
	m, _ := New([]string{
		"email/[^@]+/USER",
		"ip/\\d+\.\\d+\.\\d+\.\\d+/x.x.x.x",
	})
	input := raw(t, map[string]string{"email": "alice@example.com", "ip": "192.168.1.1"})
	out := m.Apply(input)
	var obj map[string]string
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if obj["email"] != "USER@example.com" {
		t.Errorf("email = %q", obj["email"])
	}
	if obj["ip"] != "x.x.x.x" {
		t.Errorf("ip = %q", obj["ip"])
	}
}
