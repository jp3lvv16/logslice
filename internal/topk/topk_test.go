package topk_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logslice/internal/topk"
)

func msg(m map[string]any) json.RawMessage {
	b, _ := json.Marshal(m)
	return b
}

func TestNew_EmptyField(t *testing.T) {
	_, err := topk.New("", 5)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ZeroK(t *testing.T) {
	_, err := topk.New("level", 0)
	if err == nil {
		t.Fatal("expected error for k=0")
	}
}

func TestNew_NegativeK(t *testing.T) {
	_, err := topk.New("level", -1)
	if err == nil {
		t.Fatal("expected error for negative k")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := topk.New("level", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestAdd_CountsFieldValues(t *testing.T) {
	tr, _ := topk.New("level", 10)
	for i := 0; i < 5; i++ {
		tr.Add(msg(map[string]any{"level": "error"}))
	}
	for i := 0; i < 3; i++ {
		tr.Add(msg(map[string]any{"level": "warn"}))
	}
	tr.Add(msg(map[string]any{"level": "info"}))

	top := tr.Top()
	if len(top) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(top))
	}
	if top[0].Value != "error" || top[0].Count != 5 {
		t.Errorf("expected error/5, got %v", top[0])
	}
	if top[1].Value != "warn" || top[1].Count != 3 {
		t.Errorf("expected warn/3, got %v", top[1])
	}
}

func TestTop_LimitedToK(t *testing.T) {
	tr, _ := topk.New("level", 2)
	for _, v := range []string{"a", "a", "a", "b", "b", "c"} {
		tr.Add(msg(map[string]any{"level": v}))
	}
	top := tr.Top()
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Value != "a" {
		t.Errorf("expected 'a' first, got %q", top[0].Value)
	}
}

func TestAdd_SkipsMissingField(t *testing.T) {
	tr, _ := topk.New("level", 5)
	tr.Add(msg(map[string]any{"msg": "hello"}))
	if len(tr.Top()) != 0 {
		t.Error("expected no entries for missing field")
	}
}

func TestAdd_SkipsMalformedJSON(t *testing.T) {
	tr, _ := topk.New("level", 5)
	tr.Add(json.RawMessage(`not-json`))
	if len(tr.Top()) != 0 {
		t.Error("expected no entries for malformed JSON")
	}
}

func TestReset_ClearsCounts(t *testing.T) {
	tr, _ := topk.New("level", 5)
	tr.Add(msg(map[string]any{"level": "error"}))
	tr.Reset()
	if len(tr.Top()) != 0 {
		t.Error("expected empty top after reset")
	}
}
