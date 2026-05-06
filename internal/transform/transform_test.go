package transform

import (
	"encoding/json"
	"testing"
)

func msg(t *testing.T, s string) json.RawMessage {
	t.Helper()
	return json.RawMessage(s)
}

func TestNew_ValidSpecs(t *testing.T) {
	_, err := New([]string{"level=severity", "msg=message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidSpec(t *testing.T) {
	_, err := New([]string{"noequalssign"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNew_InvalidValueRewrite(t *testing.T) {
	_, err := New([]string{"level:noeq"})
	if err == nil {
		t.Fatal("expected error for bad value-rewrite spec")
	}
}

func TestApply_Rename(t *testing.T) {
	tr, _ := New([]string{"level=severity"})
	out := tr.Apply(msg(t, `{"level":"info","msg":"hello"}`))

	var m map[string]json.RawMessage
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := m["severity"]; !ok {
		t.Error("expected key 'severity'")
	}
	if _, ok := m["level"]; ok {
		t.Error("old key 'level' should be removed")
	}
}

func TestApply_ValueRewrite(t *testing.T) {
	tr, _ := New([]string{"level:info=INFO"})
	out := tr.Apply(msg(t, `{"level":"info"}`))

	var m map[string]string
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["level"] != "INFO" {
		t.Errorf("expected INFO, got %q", m["level"])
	}
}

func TestApply_ValueRewrite_NoMatch(t *testing.T) {
	tr, _ := New([]string{"level:debug=DEBUG"})
	out := tr.Apply(msg(t, `{"level":"info"}`))

	var m map[string]string
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["level"] != "info" {
		t.Errorf("value should be unchanged, got %q", m["level"])
	}
}

func TestApply_MissingKey(t *testing.T) {
	tr, _ := New([]string{"missing=other"})
	raw := msg(t, `{"level":"info"}`)
	out := tr.Apply(raw)
	if string(out) != string(raw) {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_NonJSON_Passthrough(t *testing.T) {
	tr, _ := New([]string{"level=severity"})
	raw := msg(t, `not json at all`)
	out := tr.Apply(raw)
	if string(out) != string(raw) {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_NoRules(t *testing.T) {
	tr, _ := New(nil)
	raw := msg(t, `{"level":"info"}`)
	out := tr.Apply(raw)
	if string(out) != string(raw) {
		t.Errorf("expected identical output, got %s", out)
	}
}
