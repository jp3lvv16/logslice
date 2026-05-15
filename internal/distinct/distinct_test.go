package distinct

import (
	"encoding/json"
	"errors"
	"testing"
)

func msg(field, value string) json.RawMessage {
	b, _ := json.Marshal(map[string]string{field: value})
	return b
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", 0)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
	if !errors.Is(err, errEmptyField) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NegativeWindow_TreatedAsUnlimited(t *testing.T) {
	d, err := New("id", -5)
	if err != nil {
		t.Fatal(err)
	}
	if d.window != 0 {
		t.Fatalf("expected window 0, got %d", d.window)
	}
}

func TestNew_ValidParams(t *testing.T) {
	d, err := New("id", 10)
	if err != nil {
		t.Fatal(err)
	}
	if d.field != "id" {
		t.Fatalf("expected field 'id', got %q", d.field)
	}
}

func TestKeep_FirstOccurrence(t *testing.T) {
	d, _ := New("id", 0)
	if !d.Keep(msg("id", "abc")) {
		t.Fatal("expected first occurrence to be kept")
	}
}

func TestKeep_SecondOccurrence(t *testing.T) {
	d, _ := New("id", 0)
	d.Keep(msg("id", "abc"))
	if d.Keep(msg("id", "abc")) {
		t.Fatal("expected duplicate to be dropped")
	}
}

func TestKeep_DifferentValues(t *testing.T) {
	d, _ := New("id", 0)
	if !d.Keep(msg("id", "a")) {
		t.Fatal("expected first value kept")
	}
	if !d.Keep(msg("id", "b")) {
		t.Fatal("expected second distinct value kept")
	}
	if d.Len() != 2 {
		t.Fatalf("expected 2 tracked values, got %d", d.Len())
	}
}

func TestKeep_MalformedJSON_PassThrough(t *testing.T) {
	d, _ := New("id", 0)
	if !d.Keep(json.RawMessage(`not json`)) {
		t.Fatal("expected malformed JSON to pass through")
	}
}

func TestKeep_MissingField_PassThrough(t *testing.T) {
	d, _ := New("id", 0)
	if !d.Keep(msg("other", "x")) {
		t.Fatal("expected missing field to pass through")
	}
}

func TestKeep_WindowEviction(t *testing.T) {
	d, _ := New("id", 2)
	d.Keep(msg("id", "a")) // seen: [a]
	d.Keep(msg("id", "b")) // seen: [a, b]
	d.Keep(msg("id", "c")) // seen: [b, c] — "a" evicted

	// "a" should be accepted again after eviction
	if !d.Keep(msg("id", "a")) {
		t.Fatal("expected evicted value to be accepted again")
	}
	if d.Len() != 2 {
		t.Fatalf("expected window size 2, got %d", d.Len())
	}
}
