package flatten

import (
	"encoding/json"
	"testing"
)

func TestNew_DefaultSeparator(t *testing.T) {
	f := New("")
	if f.separator != "." {
		t.Fatalf("expected default separator '.', got %q", f.separator)
	}
}

func TestApply_FlatObject(t *testing.T) {
	f := New(".")
	input := json.RawMessage(`{"level":"info","msg":"hello"}`)
	out, err := f.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatal(err)
	}
	if m["level"] != "info" || m["msg"] != "hello" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestApply_NestedObject(t *testing.T) {
	f := New(".")
	input := json.RawMessage(`{"request":{"method":"GET","path":"/api"}}`)
	out, err := f.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatal(err)
	}
	if m["request.method"] != "GET" {
		t.Fatalf("expected request.method=GET, got %v", m)
	}
	if m["request.path"] != "/api" {
		t.Fatalf("expected request.path=/api, got %v", m)
	}
	if _, ok := m["request"]; ok {
		t.Fatal("nested key 'request' should not exist after flattening")
	}
}

func TestApply_DeeplyNested(t *testing.T) {
	f := New(".")
	input := json.RawMessage(`{"a":{"b":{"c":42}}}`)
	out, err := f.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatal(err)
	}
	val, ok := m["a.b.c"]
	if !ok {
		t.Fatal("expected key a.b.c")
	}
	if val.(float64) != 42 {
		t.Fatalf("expected 42, got %v", val)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f := New("_")
	input := json.RawMessage(`{"http":{"status":200}}`)
	out, err := f.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["http_status"]; !ok {
		t.Fatalf("expected key http_status, got keys: %v", m)
	}
}

func TestApply_NonObjectPassthrough(t *testing.T) {
	f := New(".")
	input := json.RawMessage(`"just a string"`)
	out, err := f.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `"just a string"` {
		t.Fatalf("expected passthrough, got %s", out)
	}
}
