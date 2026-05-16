package clamp_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logslice/internal/clamp"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_EmptyField(t *testing.T) {
	_, err := clamp.New("", 0, 100)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_MinExceedsMax(t *testing.T) {
	_, err := clamp.New("val", 50, 10)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNew_Valid(t *testing.T) {
	c, err := clamp.New("val", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Clamper")
	}
}

func TestApply_WithinRange_Unchanged(t *testing.T) {
	c, _ := clamp.New("val", 0, 100)
	in := raw(`{"val":50}`)
	out, err := c.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]float64
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatal(err)
	}
	if obj["val"] != 50 {
		t.Errorf("expected 50, got %g", obj["val"])
	}
}

func TestApply_BelowMin_RaisedToMin(t *testing.T) {
	c, _ := clamp.New("val", 10, 100)
	out, err := c.Apply(raw(`{"val":-5}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]float64
	json.Unmarshal(out, &obj)
	if obj["val"] != 10 {
		t.Errorf("expected 10, got %g", obj["val"])
	}
}

func TestApply_AboveMax_LoweredToMax(t *testing.T) {
	c, _ := clamp.New("val", 0, 100)
	out, err := c.Apply(raw(`{"val":999}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]float64
	json.Unmarshal(out, &obj)
	if obj["val"] != 100 {
		t.Errorf("expected 100, got %g", obj["val"])
	}
}

func TestApply_MissingField_Passthrough(t *testing.T) {
	c, _ := clamp.New("score", 0, 10)
	in := raw(`{"other":"x"}`)
	out, err := c.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(in) {
		t.Errorf("expected unchanged message")
	}
}

func TestApply_NonNumericField_Passthrough(t *testing.T) {
	c, _ := clamp.New("val", 0, 100)
	in := raw(`{"val":"hello"}`)
	out, err := c.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(in) {
		t.Errorf("expected unchanged message")
	}
}

func TestApply_InvalidJSON_ReturnsError(t *testing.T) {
	c, _ := clamp.New("val", 0, 100)
	_, err := c.Apply(raw(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
