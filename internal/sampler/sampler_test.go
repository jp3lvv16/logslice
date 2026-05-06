package sampler_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/sampler"
)

func TestNew_ValidRate(t *testing.T) {
	s, err := sampler.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 3 {
		t.Errorf("expected rate 3, got %d", s.Rate())
	}
}

func TestNew_RateOne(t *testing.T) {
	s, err := sampler.New(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 1 {
		t.Errorf("expected rate 1, got %d", s.Rate())
	}
}

func TestNew_InvalidRate(t *testing.T) {
	for _, rate := range []int{0, -1, -100} {
		_, err := sampler.New(rate)
		if err == nil {
			t.Errorf("expected error for rate %d, got nil", rate)
		}
	}
}

func TestKeep_EveryNth(t *testing.T) {
	s, _ := sampler.New(3)
	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, want := range expected {
		got := s.Keep()
		if got != want {
			t.Errorf("call %d: expected Keep()=%v, got %v", i+1, want, got)
		}
	}
}

func TestKeep_RateOne_AlwaysTrue(t *testing.T) {
	s, _ := sampler.New(1)
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Errorf("call %d: expected Keep()=true for rate 1", i+1)
		}
	}
}

func TestCount_TracksTotal(t *testing.T) {
	s, _ := sampler.New(2)
	for i := 0; i < 5; i++ {
		s.Keep()
	}
	if s.Count() != 5 {
		t.Errorf("expected count 5, got %d", s.Count())
	}
}

func TestReset_ClearsCounter(t *testing.T) {
	s, _ := sampler.New(3)
	s.Keep()
	s.Keep()
	s.Reset()
	if s.Count() != 0 {
		t.Errorf("expected count 0 after reset, got %d", s.Count())
	}
	// After reset, the cycle should restart
	if s.Keep() {
		t.Error("expected Keep()=false on first call after reset with rate 3")
	}
}
