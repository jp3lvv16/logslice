package coalesce

import (
	"encoding/json"
	"testing"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_ValidSpec(t *testing.T) {
	_, err := New("message", []string{"msg", "text"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyTarget(t *testing.T) {
	_, err := New("", []string{"msg"}, false)
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestNew_NoCandidates(t *testing.T) {
	_, err := New("message", nil, false)
	if err == nil {
		t.Fatal("expected error for empty candidates")
	}
}

func TestApply_FirstCandidateWins(t *testing.T) {
	c, _ := New("out", []string{"a", "b"}, false)
	got, err := c.Apply(raw(`{"a":"hello","b":"world"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]json.RawMessage
	json.Unmarshal(got, &obj)
	if string(obj["out"]) != `"hello"` {
		t.Errorf("expected \"hello\", got %s", obj["out"])
	}
}

func TestApply_SkipsNullAndEmpty(t *testing.T) {
	c, _ := New("out", []string{"a", "b", "c"}, false)
	got, err := c.Apply(raw(`{"a":null,"b":"","c":"found"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]json.RawMessage
	json.Unmarshal(got, &obj)
	if string(obj["out"]) != `"found"` {
		t.Errorf("expected \"found\", got %s", obj["out"])
	}
}

func TestApply_RemovesSrcWhenKeepFalse(t *testing.T) {
	c, _ := New("out", []string{"src"}, false)
	got, _ := c.Apply(raw(`{"src":"value"}`))
	var obj map[string]json.RawMessage
	json.Unmarshal(got, &obj)
	if _, ok := obj["src"]; ok {
		t.Error("expected src field to be removed")
	}
}

func TestApply_KeepsSrcWhenKeepTrue(t *testing.T) {
	c, _ := New("out", []string{"src"}, true)
	got, _ := c.Apply(raw(`{"src":"value"}`))
	var obj map[string]json.RawMessage
	json.Unmarshal(got, &obj)
	if _, ok := obj["src"]; !ok {
		t.Error("expected src field to be retained")
	}
}

func TestApply_NoCandidatePresent_Unchanged(t *testing.T) {
	c, _ := New("out", []string{"missing"}, false)
	input := raw(`{"level":"info"}`)
	got, _ := c.Apply(input)
	if string(got) != string(input) {
		t.Errorf("expected unchanged output, got %s", got)
	}
}

func TestApply_MalformedJSON_Passthrough(t *testing.T) {
	c, _ := New("out", []string{"a"}, false)
	input := raw(`not json`)
	got, _ := c.Apply(input)
	if string(got) != string(input) {
		t.Errorf("expected passthrough for malformed JSON")
	}
}
