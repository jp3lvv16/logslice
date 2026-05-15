package numeric_test

import (
	"math"
	"testing"

	"github.com/yourorg/logslice/internal/numeric"
)

func TestNew_EmptyField(t *testing.T) {
	_, err := numeric.New("", 0, 100)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_MinExceedsMax(t *testing.T) {
	_, err := numeric.New("latency", 100, 10)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNew_Valid(t *testing.T) {
	f, err := numeric.New("latency", 0, 500)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "latency" {
		t.Errorf("field = %q, want latency", f.Field())
	}
	lo, hi := f.Bounds()
	if lo != 0 || hi != 500 {
		t.Errorf("bounds = [%g, %g], want [0, 500]", lo, hi)
	}
}

func TestKeep_WithinRange(t *testing.T) {
	f, _ := numeric.New("status", 200, 299)
	line := []byte(`{"status":200,"msg":"ok"}`)
	if !f.Keep(line) {
		t.Error("expected Keep=true for value within range")
	}
}

func TestKeep_BelowMin(t *testing.T) {
	f, _ := numeric.New("status", 200, 299)
	if f.Keep([]byte(`{"status":100}`)) {
		t.Error("expected Keep=false for value below min")
	}
}

func TestKeep_AboveMax(t *testing.T) {
	f, _ := numeric.New("status", 200, 299)
	if f.Keep([]byte(`{"status":500}`)) {
		t.Error("expected Keep=false for value above max")
	}
}

func TestKeep_MissingField(t *testing.T) {
	f, _ := numeric.New("latency", 0, 100)
	if f.Keep([]byte(`{"msg":"no latency here"}`)) {
		t.Error("expected Keep=false when field is absent")
	}
}

func TestKeep_NonNumericField(t *testing.T) {
	f, _ := numeric.New("latency", 0, 100)
	if f.Keep([]byte(`{"latency":"fast"}`)) {
		t.Error("expected Keep=false for non-numeric value")
	}
}

func TestKeep_MalformedJSON(t *testing.T) {
	f, _ := numeric.New("x", 0, 1)
	if f.Keep([]byte(`not json`)) {
		t.Error("expected Keep=false for malformed JSON")
	}
}

func TestKeep_OpenUpperBound(t *testing.T) {
	f, _ := numeric.New("score", 0, math.Inf(1))
	if !f.Keep([]byte(`{"score":9999999}`)) {
		t.Error("expected Keep=true with open upper bound")
	}
}

func TestString_ContainsField(t *testing.T) {
	f, _ := numeric.New("latency", 10, 500)
	if s := f.String(); s == "" {
		t.Error("String() returned empty")
	}
}
