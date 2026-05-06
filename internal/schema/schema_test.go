package schema

import (
	"encoding/json"
	"testing"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_NoSpecs(t *testing.T) {
	v, err := New(nil)
	if err != nil || v == nil {
		t.Fatal("expected non-nil validator with no error")
	}
}

func TestNew_ValidSpecs(t *testing.T) {
	specs := []string{"level:string", "ts:number", "ok:boolean", "msg"}
	v, err := New(specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v.rules) != 4 {
		t.Fatalf("expected 4 rules, got %d", len(v.rules))
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New([]string{"field:object"})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestValidate_NoRules(t *testing.T) {
	v, _ := New(nil)
	if err := v.Validate(raw(`{"a":1}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_AllPresent(t *testing.T) {
	v, _ := New([]string{"level:string", "ts:number", "ok:boolean"})
	msg := raw(`{"level":"info","ts":1.0,"ok":true}`)
	if err := v.Validate(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_MissingField(t *testing.T) {
	v, _ := New([]string{"level:string", "ts:number"})
	if err := v.Validate(raw(`{"level":"info"}`)); err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestValidate_WrongType_String(t *testing.T) {
	v, _ := New([]string{"level:string"})
	if err := v.Validate(raw(`{"level":42}`)); err == nil {
		t.Fatal("expected type error")
	}
}

func TestValidate_WrongType_Number(t *testing.T) {
	v, _ := New([]string{"ts:number"})
	if err := v.Validate(raw(`{"ts":"not-a-number"}`)); err == nil {
		t.Fatal("expected type error")
	}
}

func TestValidate_WrongType_Boolean(t *testing.T) {
	v, _ := New([]string{"ok:boolean"})
	if err := v.Validate(raw(`{"ok":"yes"}`)); err == nil {
		t.Fatal("expected type error")
	}
}

func TestValidate_AnyType_AcceptsAll(t *testing.T) {
	v, _ := New([]string{"data"})
	for _, val := range []string{`{"data":1}`, `{"data":"x"}`, `{"data":true}`, `{"data":null}`} {
		if err := v.Validate(raw(val)); err != nil {
			t.Fatalf("unexpected error for %s: %v", val, err)
		}
	}
}

func TestValidate_InvalidJSON(t *testing.T) {
	v, _ := New([]string{"level"})
	if err := v.Validate(raw(`not-json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
