package join

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNew_EmptyKey(t *testing.T) {
	_, err := New("", time.Second)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNew_ZeroWindow(t *testing.T) {
	_, err := New("id", 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_Valid(t *testing.T) {
	j, err := New("id", time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j == nil {
		t.Fatal("expected non-nil joiner")
	}
}

func TestAdd_BuffersFirstOccurrence(t *testing.T) {
	j, _ := New("id", time.Second)
	out, err := j.Add([]byte(`{"id":"x","a":1}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil on first occurrence, got %s", out)
	}
}

func TestAdd_MergesSecondOccurrence(t *testing.T) {
	j, _ := New("id", time.Second)
	_, _ = j.Add([]byte(`{"id":"x","a":1}`))
	out, err := j.Add([]byte(`{"id":"x","b":2}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatal("expected merged output on second occurrence")
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["a"]; !ok {
		t.Error("merged output missing field 'a'")
	}
	if _, ok := m["b"]; !ok {
		t.Error("merged output missing field 'b'")
	}
}

func TestAdd_MissingKey_ReturnsNil(t *testing.T) {
	j, _ := New("id", time.Second)
	out, err := j.Add([]byte(`{"other":"value"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatal("expected nil when key absent")
	}
}

func TestAdd_MalformedJSON_ReturnsError(t *testing.T) {
	j, _ := New("id", time.Second)
	_, err := j.Add([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestFlush_ReturnsPending(t *testing.T) {
	j, _ := New("id", time.Second)
	_, _ = j.Add([]byte(`{"id":"a","v":1}`))
	_, _ = j.Add([]byte(`{"id":"b","v":2}`))
	lines, err := j.Flush()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 flushed lines, got %d", len(lines))
	}
}

func TestFlush_ClearsBuffer(t *testing.T) {
	j, _ := New("id", time.Second)
	_, _ = j.Add([]byte(`{"id":"a","v":1}`))
	_, _ = j.Flush()
	lines, _ := j.Flush()
	if len(lines) != 0 {
		t.Fatalf("expected empty flush after clear, got %d", len(lines))
	}
}
