package throttle

import (
	"encoding/json"
	"testing"
	"time"
)

var dummyMsg = json.RawMessage(`{"msg":"hello"}`)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(10, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_Valid(t *testing.T) {
	th, err := New(5, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttler")
	}
}

func TestAllow_UnlimitedWhenMaxZero(t *testing.T) {
	th, _ := New(0, time.Second)
	for i := 0; i < 1000; i++ {
		if !th.Allow(dummyMsg) {
			t.Fatal("expected all lines allowed when max=0")
		}
	}
	if th.Dropped() != 0 {
		t.Fatalf("expected 0 dropped, got %d", th.Dropped())
	}
}

func TestAllow_StopsAtMax(t *testing.T) {
	th, _ := New(3, time.Second)
	allowed := 0
	for i := 0; i < 10; i++ {
		if th.Allow(dummyMsg) {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 allowed, got %d", allowed)
	}
	if th.Dropped() != 7 {
		t.Fatalf("expected 7 dropped, got %d", th.Dropped())
	}
}

func TestAllow_RefillsAfterWindow(t *testing.T) {
	th, _ := New(2, 50*time.Millisecond)
	th.Allow(dummyMsg)
	th.Allow(dummyMsg)
	if th.Allow(dummyMsg) {
		t.Fatal("expected third call to be dropped")
	}
	time.Sleep(60 * time.Millisecond)
	if !th.Allow(dummyMsg) {
		t.Fatal("expected allow after window reset")
	}
}

func TestDropped_InitiallyZero(t *testing.T) {
	th, _ := New(5, time.Second)
	if th.Dropped() != 0 {
		t.Fatalf("expected 0 dropped initially, got %d", th.Dropped())
	}
}
