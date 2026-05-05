package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func rawMsg(t *testing.T, v interface{}) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return json.RawMessage(b)
}

func TestWrite_FormatJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatJSON, nil)

	msg := rawMsg(t, map[string]interface{}{"level": "info", "msg": "hello"})
	if err := w.Write(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
		t.Errorf("expected JSON object, got: %s", line)
	}
}

func TestWrite_FormatPretty(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatPretty, nil)

	msg := rawMsg(t, map[string]interface{}{"level": "warn"})
	if err := w.Write(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "\n") {
		t.Error("expected indented (multi-line) output")
	}
}

func TestWrite_FormatText(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatText, nil)

	msg := rawMsg(t, map[string]interface{}{"level": "error", "msg": "boom"})
	if err := w.Write(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	if !strings.Contains(line, "=") {
		t.Errorf("expected key=value format, got: %s", line)
	}
}

func TestWrite_FieldProjection(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatJSON, []string{"msg"})

	msg := rawMsg(t, map[string]interface{}{"level": "info", "msg": "hello", "ts": "2024-01-01"})
	if err := w.Write(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &result); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if _, ok := result["msg"]; !ok {
		t.Error("expected 'msg' field in output")
	}
	if _, ok := result["level"]; ok {
		t.Error("expected 'level' field to be projected out")
	}
}

func TestWrite_MalformedInput(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatJSON, nil)

	if err := w.Write(json.RawMessage(`not-json`)); err == nil {
		t.Error("expected error for malformed JSON input")
	}
}
