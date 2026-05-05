package timerange

import (
	"testing"
	"time"
)

func TestParse_ValidRange(t *testing.T) {
	tr, err := Parse("2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Start == nil || tr.End == nil {
		t.Fatal("expected both Start and End to be set")
	}
}

func TestParse_OpenStart(t *testing.T) {
	tr, err := Parse("", "2024-06-01T12:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Start != nil {
		t.Error("expected Start to be nil")
	}
	if tr.End == nil {
		t.Error("expected End to be set")
	}
}

func TestParse_EndBeforeStart(t *testing.T) {
	_, err := Parse("2024-06-02T00:00:00Z", "2024-06-01T00:00:00Z")
	if err == nil {
		t.Fatal("expected error for end before start")
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := Parse("not-a-date", "")
	if err == nil {
		t.Fatal("expected error for invalid date format")
	}
}

func TestParse_AlternativeLayouts(t *testing.T) {
	cases := []string{
		"2024-03-15T10:30:00",
		"2024-03-15 10:30:00",
		"2024-03-15",
	}
	for _, s := range cases {
		_, err := Parse(s, "")
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", s, err)
		}
	}
}

func TestContains(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	tr := TimeRange{Start: &start, End: &end}

	cases := []struct {
		t      time.Time
		wantIn bool
	}{
		{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), true},
		{start, true},
		{end, true},
		{time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), false},
	}

	for _, c := range cases {
		got := tr.Contains(c.t)
		if got != c.wantIn {
			t.Errorf("Contains(%v) = %v, want %v", c.t, got, c.wantIn)
		}
	}
}

func TestIsZero(t *testing.T) {
	var tr TimeRange
	if !tr.IsZero() {
		t.Error("expected zero TimeRange to report IsZero() = true")
	}

	tr2, _ := Parse("2024-01-01T00:00:00Z", "")
	if tr2.IsZero() {
		t.Error("expected non-zero TimeRange to report IsZero() = false")
	}
}
