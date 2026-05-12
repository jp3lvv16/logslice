package window

import (
	"encoding/json"
	"testing"
	"time"
)

func msg(s string) json.RawMessage { return json.RawMessage(s) }

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNew_NegativeSize(t *testing.T) {
	_, err := New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative size")
	}
}

func TestNew_ZeroSize(t *testing.T) {
	w, err := New(0)
	if err != nil || w == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdd_ZeroSize_EachMessageFlushed(t *testing.T) {
	w, _ := New(0)
	for i := 0; i < 3; i++ {
		g, flush := w.Add(t0, msg(`{}`))
		if !flush || g == nil || len(g.Messages) != 1 {
			t.Fatalf("iteration %d: expected immediate flush", i)
		}
	}
}

func TestAdd_SameBucket_NoFlush(t *testing.T) {
	w, _ := New(time.Minute)
	for i := 0; i < 4; i++ {
		ts := t0.Add(time.Duration(i) * 10 * time.Second)
		g, flush := w.Add(ts, msg(`{}`))
		if flush || g != nil {
			t.Fatalf("iteration %d: unexpected flush", i)
		}
	}
	if len(w.current.Messages) != 4 {
		t.Fatalf("expected 4 messages, got %d", len(w.current.Messages))
	}
}

func TestAdd_NewBucket_FlushesOld(t *testing.T) {
	w, _ := New(time.Minute)
	w.Add(t0, msg(`{"a":1}`))
	w.Add(t0.Add(30*time.Second), msg(`{"a":2}`))

	// this timestamp crosses the 1-minute boundary
	g, flush := w.Add(t0.Add(61*time.Second), msg(`{"a":3}`))
	if !flush || g == nil {
		t.Fatal("expected flush when crossing bucket boundary")
	}
	if len(g.Messages) != 2 {
		t.Fatalf("expected 2 messages in flushed group, got %d", len(g.Messages))
	}
	if len(w.current.Messages) != 1 {
		t.Fatalf("expected 1 message in new bucket, got %d", len(w.current.Messages))
	}
}

func TestFlush_DrainsFinalBucket(t *testing.T) {
	w, _ := New(time.Minute)
	w.Add(t0, msg(`{"x":1}`))
	w.Add(t0.Add(10*time.Second), msg(`{"x":2}`))

	g := w.Flush()
	if g == nil || len(g.Messages) != 2 {
		t.Fatalf("expected 2 messages from Flush, got %v", g)
	}
	if w.current != nil {
		t.Fatal("expected current to be nil after Flush")
	}
}

func TestFlush_EmptyWindow_ReturnsNil(t *testing.T) {
	w, _ := New(time.Minute)
	if g := w.Flush(); g != nil {
		t.Fatalf("expected nil from empty Flush, got %v", g)
	}
}
