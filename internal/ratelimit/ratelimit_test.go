package ratelimit

import (
	"testing"
	"time"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
	_, err = New(-5)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNew_ValidRate(t *testing.T) {
	l, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestAllow_BurstUpToMax(t *testing.T) {
	l, _ := New(5)
	// At creation tokens == max == 5; first 5 calls should succeed.
	allowed := 0
	for i := 0; i < 5; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected 5 allowed, got %d", allowed)
	}
	// Next call should be denied (no tokens left, no time elapsed).
	if l.Allow() {
		t.Fatal("expected deny after burst exhausted")
	}
}

func TestAllow_RefillOverTime(t *testing.T) {
	l, _ := New(100) // 100 lines/sec
	// Drain all tokens.
	for l.Allow() {
	}

	// Inject a fake clock that jumps forward 50ms => should add 5 tokens.
	base := time.Now()
	calls := 0
	l.now = func() time.Time {
		calls++
		return base.Add(time.Duration(calls) * 50 * time.Millisecond)
	}
	l.lastTick = base

	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed == 0 {
		t.Fatal("expected some lines allowed after time advance")
	}
}

func TestUnlimited_AlwaysAllows(t *testing.T) {
	l := Unlimited()
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatalf("unlimited limiter denied line %d", i)
		}
	}
}

func TestIsUnlimited(t *testing.T) {
	l := Unlimited()
	if !l.isUnlimited() {
		t.Fatal("expected unlimited limiter to report isUnlimited")
	}
	l2, _ := New(10)
	if l2.isUnlimited() {
		t.Fatal("expected rate-limited limiter to not report isUnlimited")
	}
}
