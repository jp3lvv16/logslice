package grep_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/grep"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := grep.New("", false)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := grep.New("[invalid", false)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidPattern(t *testing.T) {
	g, err := grep.New(`error`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Pattern() != "error" {
		t.Errorf("Pattern() = %q, want %q", g.Pattern(), "error")
	}
	if g.Inverted() {
		t.Error("Inverted() should be false")
	}
}

func TestKeep_MatchingLine(t *testing.T) {
	g, _ := grep.New(`"level":"error"`, false)
	raw := []byte(`{"level":"error","msg":"boom"}`)
	if !g.Keep(raw) {
		t.Error("expected Keep to return true for matching line")
	}
}

func TestKeep_NonMatchingLine(t *testing.T) {
	g, _ := grep.New(`"level":"error"`, false)
	raw := []byte(`{"level":"info","msg":"ok"}`)
	if g.Keep(raw) {
		t.Error("expected Keep to return false for non-matching line")
	}
}

func TestKeep_InvertMatch(t *testing.T) {
	g, _ := grep.New(`debug`, true)
	match := []byte(`{"level":"debug","msg":"verbose"}`)
	nomatch := []byte(`{"level":"info","msg":"normal"}`)

	if g.Keep(match) {
		t.Error("inverted: expected Keep to return false for matching line")
	}
	if !g.Keep(nomatch) {
		t.Error("inverted: expected Keep to return true for non-matching line")
	}
}

func TestKeep_CaseInsensitivePattern(t *testing.T) {
	g, err := grep.New(`(?i)ERROR`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw := []byte(`{"level":"error","msg":"oops"}`)
	if !g.Keep(raw) {
		t.Error("expected case-insensitive match")
	}
}

func TestKeep_EmptyLine(t *testing.T) {
	g, _ := grep.New(`error`, false)
	if g.Keep([]byte{}) {
		t.Error("expected Keep to return false for empty line")
	}
}
