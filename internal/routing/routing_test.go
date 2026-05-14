package routing

import (
	"testing"
)

func TestNew_NoSpecs(t *testing.T) {
	r, err := New(nil, "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.defaultDest != "default" {
		t.Errorf("expected default dest 'default', got %q", r.defaultDest)
	}
}

func TestNew_InvalidSpec_MissingDest(t *testing.T) {
	_, err := New([]string{"level=error"}, "default")
	if err == nil {
		t.Fatal("expected error for missing dest")
	}
}

func TestNew_InvalidSpec_MissingField(t *testing.T) {
	_, err := New([]string{"error:errors"}, "default")
	if err == nil {
		t.Fatal("expected error for missing field prefix")
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := New([]string{"level=[invalid:dest"}, "default")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestRoute_MatchesFirstRule(t *testing.T) {
	r, err := New([]string{"level=error:errors", "level=warn:warnings"}, "other")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dest := r.Route([]byte(`{"level":"error","msg":"boom"}`))
	if dest != "errors" {
		t.Errorf("expected 'errors', got %q", dest)
	}
}

func TestRoute_MatchesSecondRule(t *testing.T) {
	r, _ := New([]string{"level=error:errors", "level=warn:warnings"}, "other")
	dest := r.Route([]byte(`{"level":"warn","msg":"careful"}`))
	if dest != "warnings" {
		t.Errorf("expected 'warnings', got %q", dest)
	}
}

func TestRoute_DefaultWhenNoMatch(t *testing.T) {
	r, _ := New([]string{"level=error:errors"}, "default")
	dest := r.Route([]byte(`{"level":"info","msg":"ok"}`))
	if dest != "default" {
		t.Errorf("expected 'default', got %q", dest)
	}
}

func TestRoute_MalformedJSON(t *testing.T) {
	r, _ := New([]string{"level=error:errors"}, "default")
	dest := r.Route([]byte(`not json`))
	if dest != "default" {
		t.Errorf("expected 'default', got %q", dest)
	}
}

func TestRoute_FieldMissing(t *testing.T) {
	r, _ := New([]string{"level=error:errors"}, "default")
	dest := r.Route([]byte(`{"msg":"no level field"}`))
	if dest != "default" {
		t.Errorf("expected 'default', got %q", dest)
	}
}

func TestRoute_PartialPatternMatch(t *testing.T) {
	r, _ := New([]string{"service=auth.*:auth"}, "other")
	dest := r.Route([]byte(`{"service":"auth-service"}`))
	if dest != "auth" {
		t.Errorf("expected 'auth', got %q", dest)
	}
}
