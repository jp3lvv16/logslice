package truncate

import (
	"encoding/json"
	"testing"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_ValidParams(t *testing.T) {
	tr, err := New(80, "…")
	if err != nil || tr == nil {
		t.Fatalf("expected valid Truncator, got err=%v", err)
	}
}

func TestNew_NegativeMaxLen(t *testing.T) {
	_, err := New(-1, "")
	if err == nil {
		t.Fatal("expected error for negative maxLen")
	}
}

func TestApply_NoOp_WhenMaxLenZero(t *testing.T) {
	tr, _ := New(0, "…")
	input := raw(`{"msg":"this is a long message that should not be touched"}`)
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_TruncatesLongString(t *testing.T) {
	tr, _ := New(10, "...")
	input := raw(`{"msg":"hello world this is long"}`)
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if obj["msg"] != "hello worl..." {
		t.Errorf("unexpected value: %q", obj["msg"])
	}
}

func TestApply_PreservesShortString(t *testing.T) {
	tr, _ := New(50, "...")
	input := raw(`{"msg":"short"}`)
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	json.Unmarshal(out, &obj)
	if obj["msg"] != "short" {
		t.Errorf("expected 'short', got %q", obj["msg"])
	}
}

func TestApply_PreservesNonStringFields(t *testing.T) {
	tr, _ := New(5, "...")
	input := raw(`{"count":42,"flag":true,"ratio":3.14}`)
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]json.RawMessage
	json.Unmarshal(out, &obj)
	if string(obj["count"]) != "42" {
		t.Errorf("count changed: %s", obj["count"])
	}
}

func TestApply_PassThroughNonObject(t *testing.T) {
	tr, _ := New(10, "...")
	input := raw(`"just a string"`)
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_MultipleFields(t *testing.T) {
	tr, _ := New(4, "~")
	input := raw(`{"a":"abcdef","b":"xy","c":"123456"}`)
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	json.Unmarshal(out, &obj)
	if obj["a"] != "abcd~" {
		t.Errorf("a: got %q", obj["a"])
	}
	if obj["b"] != "xy" {
		t.Errorf("b should be unchanged, got %q", obj["b"])
	}
	if obj["c"] != "1234~" {
		t.Errorf("c: got %q", obj["c"])
	}
}
