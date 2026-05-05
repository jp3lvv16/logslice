package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestNew_ValidSpecs(t *testing.T) {
	f, err := filter.New([]string{"level=error", "service=auth.*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNew_MissingSeparator(t *testing.T) {
	_, err := filter.New([]string{"levelonly"})
	if err == nil {
		t.Fatal("expected error for missing separator")
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := filter.New([]string{"level=[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestMatch_AllRulesMatch(t *testing.T) {
	f, _ := filter.New([]string{"level=error", "service=auth"})
	entry := map[string]string{"level": "error", "service": "auth", "msg": "login failed"}
	if !f.Match(entry) {
		t.Error("expected entry to match")
	}
}

func TestMatch_PartialMatch(t *testing.T) {
	f, _ := filter.New([]string{"level=error", "service=auth"})
	entry := map[string]string{"level": "error", "service": "payments"}
	if f.Match(entry) {
		t.Error("expected entry not to match")
	}
}

func TestMatch_MissingField(t *testing.T) {
	f, _ := filter.New([]string{"level=error"})
	entry := map[string]string{"msg": "hello"}
	if f.Match(entry) {
		t.Error("expected entry not to match when field is absent")
	}
}

func TestMatch_NoRules(t *testing.T) {
	f, _ := filter.New(nil)
	entry := map[string]string{"level": "debug"}
	if !f.Match(entry) {
		t.Error("expected entry to match when there are no rules")
	}
}

func TestMatch_RegexPattern(t *testing.T) {
	f, _ := filter.New([]string{"level=^(warn|error)$"})
	if !f.Match(map[string]string{"level": "warn"}) {
		t.Error("expected 'warn' to match")
	}
	if f.Match(map[string]string{"level": "debug"}) {
		t.Error("expected 'debug' not to match")
	}
}
