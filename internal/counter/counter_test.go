package counter_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logslice/internal/counter"
)

var msg = json.RawMessage(`{"level":"info","msg":"hello"}`)

func TestNew_ZeroLimit_Unlimited(t *testing.T) {
	c := counter.New(0)
	for i := 0; i < 10_000; i++ {
		if !c.Allow(msg) {
			t.Fatalf("Allow returned false at iteration %d with unlimited counter", i)
		}
	}
}

func TestNew_NegativeLimit_TreatedAsUnlimited(t *testing.T) {
	c := counter.New(-5)
	if !c.Allow(msg) {
		t.Fatal("expected Allow to return true for negative limit")
	}
}

func TestAllow_StopsAtMax(t *testing.T) {
	const max = 3
	c := counter.New(max)

	for i := 0; i < max; i++ {
		if !c.Allow(msg) {
			t.Fatalf("Allow returned false too early at i=%d", i)
		}
	}

	if c.Allow(msg) {
		t.Fatal("Allow should return false after limit is reached")
	}
}

func TestSeen_TracksCount(t *testing.T) {
	c := counter.New(5)
	for i := 0; i < 3; i++ {
		c.Allow(msg)
	}
	if c.Seen() != 3 {
		t.Fatalf("expected Seen()=3, got %d", c.Seen())
	}
}

func TestSeen_DoesNotExceedMax(t *testing.T) {
	const max = 2
	c := counter.New(max)
	for i := 0; i < 10; i++ {
		c.Allow(msg)
	}
	if c.Seen() != max {
		t.Fatalf("expected Seen()=%d, got %d", max, c.Seen())
	}
}

func TestDone_FalseBeforeLimit(t *testing.T) {
	c := counter.New(2)
	c.Allow(msg)
	if c.Done() {
		t.Fatal("Done should be false before limit is reached")
	}
}

func TestDone_TrueAfterLimit(t *testing.T) {
	c := counter.New(1)
	c.Allow(msg)
	if !c.Done() {
		t.Fatal("Done should be true once limit is reached")
	}
}

func TestDone_AlwaysFalseWhenUnlimited(t *testing.T) {
	c := counter.New(0)
	for i := 0; i < 100; i++ {
		c.Allow(msg)
	}
	if c.Done() {
		t.Fatal("Done should always be false for unlimited counter")
	}
}
