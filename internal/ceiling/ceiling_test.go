package ceiling_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/logslice/logslice/internal/ceiling"
)

func msg(field, value string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{%q:%q}`, field, value))
}

func TestNew_EmptyField(t *testing.T) {
	_, err := ceiling.New("", 3)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ZeroMax(t *testing.T) {
	_, err := ceiling.New("level", 0)
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNew_NegativeMax(t *testing.T) {
	_, err := ceiling.New("level", -1)
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestNew_Valid(t *testing.T) {
	c, err := ceiling.New("level", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Ceiling")
	}
}

func TestAllow_StopsAtMax(t *testing.T) {
	c, _ := ceiling.New("level", 2)

	if !c.Allow(msg("level", "error")) {
		t.Error("first occurrence should be allowed")
	}
	if !c.Allow(msg("level", "error")) {
		t.Error("second occurrence should be allowed")
	}
	if c.Allow(msg("level", "error")) {
		t.Error("third occurrence should be denied")
	}
}

func TestAllow_IndependentPerValue(t *testing.T) {
	c, _ := ceiling.New("level", 1)

	if !c.Allow(msg("level", "info")) {
		t.Error("info first occurrence should be allowed")
	}
	if !c.Allow(msg("level", "warn")) {
		t.Error("warn first occurrence should be allowed")
	}
	if c.Allow(msg("level", "info")) {
		t.Error("info second occurrence should be denied")
	}
}

func TestAllow_MissingField_AlwaysAllowed(t *testing.T) {
	c, _ := ceiling.New("level", 1)

	for i := 0; i < 5; i++ {
		if !c.Allow(json.RawMessage(`{"msg":"hello"}`)) {
			t.Errorf("iteration %d: message without field should always be allowed", i)
		}
	}
}

func TestAllow_MalformedJSON_AlwaysAllowed(t *testing.T) {
	c, _ := ceiling.New("level", 1)
	if !c.Allow(json.RawMessage(`not-json`)) {
		t.Error("malformed JSON should be forwarded")
	}
}

func TestSeen_TracksCount(t *testing.T) {
	c, _ := ceiling.New("level", 5)

	for i := 0; i < 3; i++ {
		c.Allow(msg("level", "debug"))
	}
	if got := c.Seen(`"debug"`); got != 3 {
		t.Errorf("Seen = %d, want 3", got)
	}
}
