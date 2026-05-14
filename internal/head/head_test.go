package head

import (
	"testing"
)

func TestNew_NegativeLimit(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative limit, got nil")
	}
}

func TestNew_ZeroLimit_Unlimited(t *testing.T) {
	l, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 1000; i++ {
		if !l.Keep() {
			t.Fatalf("Keep() returned false at iteration %d with unlimited limiter", i)
		}
	}
	if l.Done() {
		t.Fatal("Done() should never be true for unlimited limiter")
	}
}

func TestNew_PositiveLimit(t *testing.T) {
	l, err := New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
}

func TestKeep_StopsAtMax(t *testing.T) {
	const max = 5
	l, _ := New(max)

	for i := 0; i < max; i++ {
		if !l.Keep() {
			t.Fatalf("Keep() returned false at iteration %d, expected true", i)
		}
	}
	if l.Keep() {
		t.Fatal("Keep() should return false after limit is reached")
	}
}

func TestSeen_TracksCount(t *testing.T) {
	l, _ := New(10)
	for i := 0; i < 4; i++ {
		l.Keep()
	}
	if l.Seen() != 4 {
		t.Fatalf("expected Seen() == 4, got %d", l.Seen())
	}
}

func TestSeen_DoesNotExceedMax(t *testing.T) {
	const max = 3
	l, _ := New(max)
	for i := 0; i < max+10; i++ {
		l.Keep()
	}
	if l.Seen() != max {
		t.Fatalf("expected Seen() == %d, got %d", max, l.Seen())
	}
}

func TestDone_FalseBeforeLimit(t *testing.T) {
	l, _ := New(2)
	l.Keep()
	if l.Done() {
		t.Fatal("Done() should be false before limit is reached")
	}
}

func TestDone_TrueAtLimit(t *testing.T) {
	l, _ := New(2)
	l.Keep()
	l.Keep()
	if !l.Done() {
		t.Fatal("Done() should be true once limit is reached")
	}
}
