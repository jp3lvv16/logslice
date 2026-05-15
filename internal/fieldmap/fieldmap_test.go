package fieldmap_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logslice/internal/fieldmap"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_ValidSpecs(t *testing.T) {
	m, err := fieldmap.New([]string{"level:severity", "msg:message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Mapper")
	}
}

func TestNew_EmptySpecs(t *testing.T) {
	m, err := fieldmap.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Mapper")
	}
}

func TestNew_MissingSeparator(t *testing.T) {
	_, err := fieldmap.New([]string{"levelSEVERITY"})
	if err == nil {
		t.Fatal("expected error for missing separator")
	}
}

func TestNew_EmptyOldKey(t *testing.T) {
	_, err := fieldmap.New([]string{":new"})
	if err == nil {
		t.Fatal("expected error for empty old key")
	}
}

func TestNew_EmptyNewKey(t *testing.T) {
	_, err := fieldmap.New([]string{"old:"})
	if err == nil {
		t.Fatal("expected error for empty new key")
	}
}

func TestNew_SameKey(t *testing.T) {
	_, err := fieldmap.New([]string{"level:level"})
	if err == nil {
		t.Fatal("expected error when old == new")
	}
}

func TestApply_RenamesField(t *testing.T) {
	m, _ := fieldmap.New([]string{"level:severity"})
	out := m.Apply(raw(`{"level":"info","msg":"hello"}`))

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := obj["level"]; ok {
		t.Error("old key 'level' should have been removed")
	}
	if _, ok := obj["severity"]; !ok {
		t.Error("new key 'severity' should be present")
	}
}

func TestApply_MissingSourceField(t *testing.T) {
	m, _ := fieldmap.New([]string{"level:severity"})
	out := m.Apply(raw(`{"msg":"hello"}`))

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := obj["severity"]; ok {
		t.Error("severity should not appear when source field is absent")
	}
}

func TestApply_NoOp_WhenNoRules(t *testing.T) {
	m, _ := fieldmap.New(nil)
	input := raw(`{"level":"info"}`)
	out := m.Apply(input)
	if string(out) != string(input) {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_MalformedJSON(t *testing.T) {
	m, _ := fieldmap.New([]string{"level:severity"})
	input := raw(`not-json`)
	out := m.Apply(input)
	if string(out) != string(input) {
		t.Errorf("expected original message on malformed JSON, got %s", out)
	}
}
