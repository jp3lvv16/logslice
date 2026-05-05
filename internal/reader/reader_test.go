package reader

import (
	"strings"
	"testing"
	"time"
)

func lines(input string) []Line {
	var out []Line
	for l := range Read(strings.NewReader(input)) {
		out = append(out, l)
	}
	return out
}

func TestRead_ValidJSON(t *testing.T) {
	input := `{"time":"2024-01-15T10:00:00Z","level":"info","msg":"started"}
{"time":"2024-01-15T10:01:00Z","level":"error","msg":"oops"}
`
	got := lines(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	if got[0].Fields["msg"] != "started" {
		t.Errorf("unexpected msg: %v", got[0].Fields["msg"])
	}
}

func TestRead_SkipsMalformedLines(t *testing.T) {
	input := `not json
{"ts":"2024-01-15T10:00:00Z","msg":"ok"}
also not json
`
	got := lines(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got))
	}
}

func TestRead_SkipsBlankLines(t *testing.T) {
	input := "\n{\"ts\":\"2024-01-15T10:00:00Z\",\"msg\":\"hi\"}\n\n"
	got := lines(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got))
	}
}

func TestExtractTimestamp_RFC3339(t *testing.T) {
	fields := map[string]interface{}{"time": "2024-06-01T12:30:00Z"}
	ts := extractTimestamp(fields)
	expected := time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC)
	if !ts.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, ts)
	}
}

func TestExtractTimestamp_UnixFloat(t *testing.T) {
	fields := map[string]interface{}{"ts": float64(1700000000)}
	ts := extractTimestamp(fields)
	if ts.IsZero() {
		t.Error("expected non-zero timestamp for unix float")
	}
}

func TestExtractTimestamp_Missing(t *testing.T) {
	fields := map[string]interface{}{"msg": "no time here"}
	ts := extractTimestamp(fields)
	if !ts.IsZero() {
		t.Errorf("expected zero time, got %v", ts)
	}
}

func TestExtractTimestamp_FallbackKeys(t *testing.T) {
	for _, key := range TimestampKeys {
		fields := map[string]interface{}{key: "2024-03-10T08:00:00Z"}
		ts := extractTimestamp(fields)
		if ts.IsZero() {
			t.Errorf("key %q: expected non-zero timestamp", key)
		}
	}
}
