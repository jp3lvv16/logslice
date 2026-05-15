package watermark

import (
	"encoding/json"
	"testing"
	"time"
)

func msg(ts string) []byte {
	m := map[string]string{"time": ts, "msg": "hello"}
	b, _ := json.Marshal(m)
	return b
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ValidField(t *testing.T) {
	w, err := New("time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watermark")
	}
}

func TestRecord_SkipsMalformedJSON(t *testing.T) {
	w, _ := New("time")
	w.Record([]byte("not json"))
	if w.Count() != 0 {
		t.Fatalf("expected count 0, got %d", w.Count())
	}
}

func TestRecord_SkipsMissingField(t *testing.T) {
	w, _ := New("time")
	w.Record([]byte(`{"msg":"no timestamp here"}`))
	if w.Count() != 0 {
		t.Fatalf("expected count 0, got %d", w.Count())
	}
}

func TestRecord_SkipsInvalidTimestamp(t *testing.T) {
	w, _ := New("time")
	w.Record([]byte(`{"time":"not-a-time"}`))
	if w.Count() != 0 {
		t.Fatalf("expected count 0, got %d", w.Count())
	}
}

func TestRecord_TracksMinMax(t *testing.T) {
	w, _ := New("time")
	times := []string{
		"2024-01-01T10:00:00Z",
		"2024-01-01T08:00:00Z",
		"2024-01-01T12:00:00Z",
	}
	for _, ts := range times {
		w.Record(msg(ts))
	}
	if w.Count() != 3 {
		t.Fatalf("expected count 3, got %d", w.Count())
	}
	min, ok := w.Min()
	if !ok {
		t.Fatal("expected min to be set")
	}
	if min.Hour() != 8 {
		t.Errorf("expected min hour 8, got %d", min.Hour())
	}
	max, ok := w.Max()
	if !ok {
		t.Fatal("expected max to be set")
	}
	if max.Hour() != 12 {
		t.Errorf("expected max hour 12, got %d", max.Hour())
	}
}

func TestSpan_EmptyReturnsZero(t *testing.T) {
	w, _ := New("time")
	if w.Span() != 0 {
		t.Fatalf("expected span 0, got %v", w.Span())
	}
}

func TestSpan_CorrectDuration(t *testing.T) {
	w, _ := New("time")
	w.Record(msg("2024-01-01T10:00:00Z"))
	w.Record(msg("2024-01-01T12:00:00Z"))
	expected := 2 * time.Hour
	if w.Span() != expected {
		t.Errorf("expected span %v, got %v", expected, w.Span())
	}
}

func TestMin_Max_BeforeAnyRecord(t *testing.T) {
	w, _ := New("time")
	_, okMin := w.Min()
	_, okMax := w.Max()
	if okMin || okMax {
		t.Fatal("expected false before any records")
	}
}
