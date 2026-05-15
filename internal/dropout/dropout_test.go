package dropout

import (
	"testing"
)

func TestNew_ZeroRate_Valid(t *testing.T) {
	d, err := New(0.0, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Rate() != 0.0 {
		t.Errorf("expected rate 0.0, got %v", d.Rate())
	}
}

func TestNew_ValidRate(t *testing.T) {
	d, err := New(0.5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Rate() != 0.5 {
		t.Errorf("expected rate 0.5, got %v", d.Rate())
	}
}

func TestNew_RateOne_Invalid(t *testing.T) {
	_, err := New(1.0, 0)
	if err == nil {
		t.Fatal("expected error for rate=1.0, got nil")
	}
}

func TestNew_NegativeRate_Invalid(t *testing.T) {
	_, err := New(-0.1, 0)
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
}

func TestNew_RateAboveOne_Invalid(t *testing.T) {
	_, err := New(1.5, 0)
	if err == nil {
		t.Fatal("expected error for rate > 1.0, got nil")
	}
}

func TestKeep_ZeroRate_AlwaysKeeps(t *testing.T) {
	d, _ := New(0.0, 99)
	for i := 0; i < 100; i++ {
		if !d.Keep([]byte(`{"msg":"hello"}`)) {
			t.Fatal("expected Keep to always return true for rate=0.0")
		}
	}
}

func TestKeep_HighRate_DropsApproximately(t *testing.T) {
	// With rate=0.9 and a fixed seed we expect ~90% to be dropped.
	d, _ := New(0.9, 7)
	kept := 0
	const n = 10_000
	for i := 0; i < n; i++ {
		if d.Keep([]byte(`{}`)) {
			kept++
		}
	}
	ratio := float64(kept) / float64(n)
	// Allow ±5% tolerance around the expected 10% keep rate.
	if ratio < 0.05 || ratio > 0.15 {
		t.Errorf("keep ratio %v out of expected range [0.05, 0.15]", ratio)
	}
}

func TestKeep_LowRate_KeepsMost(t *testing.T) {
	d, _ := New(0.1, 3)
	kept := 0
	const n = 10_000
	for i := 0; i < n; i++ {
		if d.Keep([]byte(`{}`)) {
			kept++
		}
	}
	ratio := float64(kept) / float64(n)
	if ratio < 0.85 || ratio > 0.95 {
		t.Errorf("keep ratio %v out of expected range [0.85, 0.95]", ratio)
	}
}
