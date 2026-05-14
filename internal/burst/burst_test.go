package burst

import (
	"testing"
	"time"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0, 5)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_Valid(t *testing.T) {
	d, err := New(time.Second, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestRecord_BelowThreshold(t *testing.T) {
	d, _ := New(time.Second, 5)
	for i := 0; i < 5; i++ {
		if d.Record() {
			t.Fatalf("burst triggered too early at event %d", i+1)
		}
	}
}

func TestRecord_ExceedsThreshold(t *testing.T) {
	d, _ := New(time.Second, 3)
	for i := 0; i < 3; i++ {
		d.Record()
	}
	if !d.Record() {
		t.Fatal("expected burst to be detected on 4th event")
	}
}

func TestRecord_EvictsExpiredEvents(t *testing.T) {
	now := time.Unix(1_000, 0)
	d, _ := New(time.Second, 2)
	d.now = func() time.Time { return now }

	// record 3 events in the past (outside window)
	past := now.Add(-2 * time.Second)
	d.now = func() time.Time { return past }
	for i := 0; i < 3; i++ {
		d.Record()
	}

	// advance time so past events are outside the window
	d.now = func() time.Time { return now }

	// these two should be within threshold=2, not burst
	if d.Record() {
		t.Fatal("should not burst: old events should be evicted")
	}
	if d.Record() {
		t.Fatal("should not burst on second fresh event")
	}
	// third fresh event exceeds threshold
	if !d.Record() {
		t.Fatal("expected burst on third fresh event")
	}
}

func TestCount_ReflectsWindow(t *testing.T) {
	now := time.Unix(2_000, 0)
	d, _ := New(time.Second, 10)
	d.now = func() time.Time { return now }

	for i := 0; i < 4; i++ {
		d.Record()
	}
	if c := d.Count(); c != 4 {
		t.Fatalf("expected count 4, got %d", c)
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	d, _ := New(time.Second, 2)
	for i := 0; i < 5; i++ {
		d.Record()
	}
	d.Reset()
	if c := d.Count(); c != 0 {
		t.Fatalf("expected 0 after reset, got %d", c)
	}
}
