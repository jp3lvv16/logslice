package sample

import (
	"encoding/json"
	"testing"
)

func msg(v int) json.RawMessage {
	return json.RawMessage([]byte(`{"n":` + string(rune('0'+v)) + `}`))
}

func TestNew_ZeroSize(t *testing.T) {
	_, err := New(0, 0)
	if err == nil {
		t.Fatal("expected error for size=0")
	}
}

func TestNew_NegativeSize(t *testing.T) {
	_, err := New(-1, 0)
	if err == nil {
		t.Fatal("expected error for size=-1")
	}
}

func TestNew_ValidSize(t *testing.T) {
	s, err := New(5, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 0 {
		t.Fatalf("expected empty reservoir, got %d", s.Len())
	}
}

func TestAdd_FillsReservoir(t *testing.T) {
	s, _ := New(3, 1)
	for i := 0; i < 3; i++ {
		s.Add(msg(i))
	}
	if s.Len() != 3 {
		t.Fatalf("expected 3, got %d", s.Len())
	}
	if s.Seen() != 3 {
		t.Fatalf("expected seen=3, got %d", s.Seen())
	}
}

func TestAdd_DoesNotExceedSize(t *testing.T) {
	s, _ := New(3, 1)
	for i := 0; i < 100; i++ {
		s.Add(msg(i % 10))
	}
	if s.Len() != 3 {
		t.Fatalf("reservoir grew beyond size: %d", s.Len())
	}
	if s.Seen() != 100 {
		t.Fatalf("expected seen=100, got %d", s.Seen())
	}
}

func TestResults_ReturnsCopy(t *testing.T) {
	s, _ := New(2, 7)
	s.Add(msg(1))
	s.Add(msg(2))
	r1 := s.Results()
	r1[0] = msg(9)
	r2 := s.Results()
	if string(r2[0]) == string(r1[0]) {
		t.Fatal("Results should return an independent copy")
	}
}

func TestAdd_SizeOne_AlwaysHoldsLatestOrEarlier(t *testing.T) {
	// With size=1 and deterministic seed we just verify reservoir stays len 1.
	s, _ := New(1, 99)
	for i := 0; i < 50; i++ {
		s.Add(msg(i % 10))
	}
	if s.Len() != 1 {
		t.Fatalf("expected len 1, got %d", s.Len())
	}
}
