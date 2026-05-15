package clip

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(1024, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_Valid(t *testing.T) {
	c, err := New(1024, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Clipper")
	}
}

func TestAllow_UnlimitedWhenMaxZero(t *testing.T) {
	c, _ := New(0, time.Second)
	msg := json.RawMessage(`{"level":"info","msg":"hello"}`)
	for i := 0; i < 100; i++ {
		if !c.Allow(msg) {
			t.Fatalf("expected Allow to return true when maxBytes=0 (iteration %d)", i)
		}
	}
}

func TestAllow_PassesWithinBudget(t *testing.T) {
	msg := json.RawMessage(`{"x":"y"}`)
	c, _ := New(int64(len(msg)*3), time.Minute)

	for i := 0; i < 3; i++ {
		if !c.Allow(msg) {
			t.Fatalf("expected Allow true on iteration %d", i)
		}
	}
}

func TestAllow_DropsWhenBudgetExceeded(t *testing.T) {
	msg := json.RawMessage(`{"x":"y"}`)
	size := int64(len(msg))
	c, _ := New(size*2, time.Minute)

	if !c.Allow(msg) {
		t.Fatal("first message should pass")
	}
	if !c.Allow(msg) {
		t.Fatal("second message should pass")
	}
	if c.Allow(msg) {
		t.Fatal("third message should be dropped")
	}
}

func TestUsed_ReflectsBytesConsumed(t *testing.T) {
	msg := json.RawMessage(`{"k":"v"}`)
	c, _ := New(1024, time.Minute)

	if c.Used() != 0 {
		t.Fatalf("expected 0 used bytes initially, got %d", c.Used())
	}
	c.Allow(msg)
	if got := c.Used(); got != int64(len(msg)) {
		t.Fatalf("expected %d used bytes, got %d", len(msg), got)
	}
}

func TestAllow_EvictsStaleEntries(t *testing.T) {
	msg := json.RawMessage(`{"a":"b"}`)
	size := int64(len(msg))
	// Very short window so we can simulate expiry.
	c, _ := New(size, 10*time.Millisecond)

	if !c.Allow(msg) {
		t.Fatal("first message should pass")
	}
	if c.Allow(msg) {
		t.Fatal("second message should be dropped (budget full)")
	}

	// Wait for the window to expire.
	time.Sleep(20 * time.Millisecond)

	if !c.Allow(msg) {
		t.Fatal("message should pass after window expires")
	}
}
