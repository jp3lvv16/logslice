package context

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var baseMsg = json.RawMessage(`{"level":"info","msg":"hello"}`)

func TestEnrich_AddsMetaFields(t *testing.T) {
	m := Meta{Line: 42, Source: "app.log", IngestedAt: time.Unix(0, 0).UTC()}
	out, err := Enrich(baseMsg, m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := obj["_line"]; !ok {
		t.Error("expected _line field")
	}
	if _, ok := obj["_source"]; !ok {
		t.Error("expected _source field")
	}
	if _, ok := obj["_ingested_at"]; !ok {
		t.Error("expected _ingested_at field")
	}
}

func TestEnrich_DoesNotOverwriteExisting(t *testing.T) {
	msg := json.RawMessage(`{"_line":99}`)
	m := Meta{Line: 1, IngestedAt: time.Now()}
	out, err := Enrich(msg, m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v, _ := obj["_line"].(float64); int(v) != 99 {
		t.Errorf("expected _line=99, got %v", obj["_line"])
	}
}

func TestEnrich_NoSource_OmitsField(t *testing.T) {
	m := Meta{Line: 1, IngestedAt: time.Now()}
	out, err := Enrich(baseMsg, m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(out), "_source") {
		t.Error("_source should be omitted when empty")
	}
}

func TestStrip_RemovesMetaFields(t *testing.T) {
	msg := json.RawMessage(`{"level":"info","_line":1,"_source":"f","_ingested_at":"2024-01-01T00:00:00Z"}`)
	out, err := Strip(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, f := range []string{"_line", "_source", "_ingested_at"} {
		if strings.Contains(string(out), f) {
			t.Errorf("field %q should have been stripped", f)
		}
	}
}

func TestStrip_PreservesOtherFields(t *testing.T) {
	msg := json.RawMessage(`{"level":"info","_line":1}`)
	out, err := Strip(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "level") {
		t.Error("level field should be preserved")
	}
}

func TestEnrich_InvalidJSON(t *testing.T) {
	_, err := Enrich(json.RawMessage(`not-json`), Meta{Line: 1, IngestedAt: time.Now()})
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
