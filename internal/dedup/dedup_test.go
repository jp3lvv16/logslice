package dedup

import (
	"fmt"
	"testing"
)

func TestNew_WindowZero_NeverDeduplicates(t *testing.T) {
	d := New(0)
	line := []byte(`{"msg":"hello"}`)
	for i := 0; i < 5; i++ {
		if d.IsDuplicate(line) {
			t.Fatalf("window=0 should never deduplicate, iteration %d", i)
		}
	}
}

func TestIsDuplicate_FirstOccurrence(t *testing.T) {
	d := New(10)
	line := []byte(`{"level":"info","msg":"startup"}`)
	if d.IsDuplicate(line) {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SecondOccurrence(t *testing.T) {
	d := New(10)
	line := []byte(`{"level":"info","msg":"startup"}`)
	d.IsDuplicate(line)
	if !d.IsDuplicate(line) {
		t.Fatal("second occurrence within window should be a duplicate")
	}
}

func TestIsDuplicate_KeyOrderIgnored(t *testing.T) {
	d := New(10)
	a := []byte(`{"msg":"hello","level":"info"}`)
	b := []byte(`{"level":"info","msg":"hello"}`)
	d.IsDuplicate(a)
	if !d.IsDuplicate(b) {
		t.Fatal("lines with same fields in different order should be duplicates")
	}
}

func TestIsDuplicate_WindowEviction(t *testing.T) {
	window := 3
	d := New(window)

	original := []byte(`{"msg":"first"}`)
	d.IsDuplicate(original) // add to window

	// Fill window with distinct lines to evict original.
	for i := 0; i < window; i++ {
		d.IsDuplicate([]byte(fmt.Sprintf(`{"msg":"filler-%d"}`, i)))
	}

	// original should have been evicted; no longer a duplicate.
	if d.IsDuplicate(original) {
		t.Fatal("line evicted from window should not be considered a duplicate")
	}
}

func TestIsDuplicate_NonJSONLine(t *testing.T) {
	d := New(10)
	line := []byte("plain text log line")
	if d.IsDuplicate(line) {
		t.Fatal("first non-JSON line should not be a duplicate")
	}
	if !d.IsDuplicate(line) {
		t.Fatal("repeated non-JSON line should be a duplicate")
	}
}

func TestIsDuplicate_DistinctLines(t *testing.T) {
	d := New(10)
	lines := [][]byte{
		[]byte(`{"msg":"alpha"}`),
		[]byte(`{"msg":"beta"}`),
		[]byte(`{"msg":"gamma"}`),
	}
	for _, l := range lines {
		if d.IsDuplicate(l) {
			t.Fatalf("distinct line %s should not be a duplicate", l)
		}
	}
}
