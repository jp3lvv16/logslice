package pivot

import (
	"encoding/json"
	"testing"
)

func TestNew_EmptyGroupKey(t *testing.T) {
	_, err := New("", "level")
	if err == nil {
		t.Fatal("expected error for empty groupKey")
	}
}

func TestNew_EmptyValueKey(t *testing.T) {
	_, err := New("service", "")
	if err == nil {
		t.Fatal("expected error for empty valueKey")
	}
}

func TestNew_Valid(t *testing.T) {
	p, err := New("service", "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Pivotter")
	}
}

func TestAdd_CountsCorrectly(t *testing.T) {
	p, _ := New("service", "level")

	lines := []string{
		`{"service":"auth","level":"info"}`,
		`{"service":"auth","level":"error"}`,
		`{"service":"auth","level":"info"}`,
		`{"service":"api","level":"info"}`,
	}
	for _, l := range lines {
		p.Add([]byte(l))
	}

	result, err := p.Result()
	if err != nil {
		t.Fatalf("Result error: %v", err)
	}

	var table map[string]map[string]int
	if err := json.Unmarshal(result, &table); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if table["auth"]["info"] != 2 {
		t.Errorf("auth/info: want 2, got %d", table["auth"]["info"])
	}
	if table["auth"]["error"] != 1 {
		t.Errorf("auth/error: want 1, got %d", table["auth"]["error"])
	}
	if table["api"]["info"] != 1 {
		t.Errorf("api/info: want 1, got %d", table["api"]["info"])
	}
}

func TestAdd_SkipsMalformedJSON(t *testing.T) {
	p, _ := New("service", "level")
	p.Add([]byte(`not json`))
	res, _ := p.Result()
	var table map[string]map[string]int
	json.Unmarshal(res, &table)
	if len(table) != 0 {
		t.Errorf("expected empty table, got %v", table)
	}
}

func TestAdd_SkipsMissingKey(t *testing.T) {
	p, _ := New("service", "level")
	p.Add([]byte(`{"service":"auth"}`))
	res, _ := p.Result()
	var table map[string]map[string]int
	json.Unmarshal(res, &table)
	if len(table) != 0 {
		t.Errorf("expected empty table for missing valueKey, got %v", table)
	}
}

func TestReset_ClearsTable(t *testing.T) {
	p, _ := New("service", "level")
	p.Add([]byte(`{"service":"auth","level":"info"}`))
	p.Reset()
	res, _ := p.Result()
	var table map[string]map[string]int
	json.Unmarshal(res, &table)
	if len(table) != 0 {
		t.Errorf("expected empty table after Reset, got %v", table)
	}
}
