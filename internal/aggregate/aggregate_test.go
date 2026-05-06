package aggregate

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func msg(s string) json.RawMessage { return json.RawMessage(s) }

func TestNew_EmptyCounter(t *testing.T) {
	c := New("level")
	if len(c.Results()) != 0 {
		t.Fatal("expected no results on a fresh counter")
	}
}

func TestAdd_CountsFieldValues(t *testing.T) {
	c := New("level")
	c.Add(msg(`{"level":"info","msg":"a"}`))
	c.Add(msg(`{"level":"info","msg":"b"}`))
	c.Add(msg(`{"level":"error","msg":"c"}`))

	res := c.Results()
	if len(res) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res))
	}
	if res[0].Value != "info" || res[0].Count != 2 {
		t.Errorf("first entry should be info:2, got %+v", res[0])
	}
	if res[1].Value != "error" || res[1].Count != 1 {
		t.Errorf("second entry should be error:1, got %+v", res[1])
	}
}

func TestAdd_MissingField(t *testing.T) {
	c := New("level")
	c.Add(msg(`{"msg":"no level here"}`))

	res := c.Results()
	if len(res) != 1 || res[0].Value != "<missing>" {
		t.Fatalf("expected <missing> entry, got %+v", res)
	}
}

func TestAdd_InvalidJSON(t *testing.T) {
	c := New("level")
	c.Add(msg(`not json`))

	res := c.Results()
	if len(res) != 1 || res[0].Value != "<missing>" {
		t.Fatalf("expected <missing> for invalid JSON, got %+v", res)
	}
}

func TestAdd_NonStringValue(t *testing.T) {
	c := New("code")
	c.Add(msg(`{"code":404}`))
	c.Add(msg(`{"code":404}`))
	c.Add(msg(`{"code":200}`))

	res := c.Results()
	if res[0].Value != "404" || res[0].Count != 2 {
		t.Errorf("expected 404:2, got %+v", res[0])
	}
}

func TestResults_SortedByCountThenValue(t *testing.T) {
	c := New("svc")
	for _, v := range []string{"b", "a", "b", "c", "a", "b"} {
		c.Add(msg(`{"svc":"` + v + `"}`))
	}
	res := c.Results()
	order := []string{"b", "a", "c"}
	for i, want := range order {
		if res[i].Value != want {
			t.Errorf("position %d: want %s got %s", i, want, res[i].Value)
		}
	}
}

func TestPrint_ContainsValues(t *testing.T) {
	c := New("level")
	c.Add(msg(`{"level":"warn"}`))
	c.Add(msg(`{"level":"warn"}`))
	c.Add(msg(`{"level":"info"}`))

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "warn") || !strings.Contains(out, "info") {
		t.Errorf("Print output missing expected values: %q", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("Print output missing count 2: %q", out)
	}
}
