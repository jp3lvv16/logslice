package template_test

import (
	"encoding/json"
	"testing"

	"github.com/logslice/logslice/internal/template"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_EmptySpec(t *testing.T) {
	_, err := template.New("")
	if err == nil {
		t.Fatal("expected error for empty spec")
	}
}

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := template.New("{{.level")
	if err == nil {
		t.Fatal("expected error for unclosed action")
	}
}

func TestNew_ValidTemplate(t *testing.T) {
	_, err := template.New("{{.level}} {{.msg}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_BasicInterpolation(t *testing.T) {
	r, _ := template.New("{{.level}}: {{.msg}}")
	out, err := r.Apply(raw(`{"level":"info","msg":"hello world"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "info: hello world"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestApply_MissingFieldRendersEmpty(t *testing.T) {
	r, _ := template.New("{{.level}} {{.missing}}")
	out, err := r.Apply(raw(`{"level":"warn"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "warn "
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	r, _ := template.New("{{.level}}")
	_, err := r.Apply(raw(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_LiteralText(t *testing.T) {
	r, _ := template.New("hello logslice")
	out, err := r.Apply(raw(`{"level":"debug"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello logslice" {
		t.Errorf("got %q", out)
	}
}

func TestApply_NumericField(t *testing.T) {
	r, _ := template.New("latency={{.latency}}ms")
	out, err := r.Apply(raw(`{"latency":42}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "latency=42ms"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}
