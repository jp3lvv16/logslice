package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/stats"
)

func TestNew_ZeroValues(t *testing.T) {
	c := stats.New()
	if c.Total != 0 || c.Matched != 0 || c.Skipped != 0 || c.Malformed != 0 {
		t.Fatalf("expected all zero counts, got %+v", c)
	}
}

func TestIncrements(t *testing.T) {
	c := stats.New()
	c.IncTotal()
	c.IncTotal()
	c.IncMatched()
	c.IncSkipped()
	c.IncSkipped()
	c.IncSkipped()
	c.IncMalformed()

	if c.Total != 2 {
		t.Errorf("Total: want 2, got %d", c.Total)
	}
	if c.Matched != 1 {
		t.Errorf("Matched: want 1, got %d", c.Matched)
	}
	if c.Skipped != 3 {
		t.Errorf("Skipped: want 3, got %d", c.Skipped)
	}
	if c.Malformed != 1 {
		t.Errorf("Malformed: want 1, got %d", c.Malformed)
	}
}

func TestElapsed_Positive(t *testing.T) {
	c := stats.New()
	if c.Elapsed() < 0 {
		t.Error("elapsed should be non-negative")
	}
}

func TestPrint_ContainsFields(t *testing.T) {
	c := stats.New()
	c.IncTotal()
	c.IncTotal()
	c.IncMatched()
	c.IncSkipped()
	c.IncMalformed()

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	for _, want := range []string{"total:", "matched:", "skipped:", "malformed:", "elapsed:"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q\ngot:\n%s", want, out)
		}
	}
}
