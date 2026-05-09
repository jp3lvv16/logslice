// Package pivot provides a Pivotter that groups JSON log lines by a key
// field and counts occurrences of each unique value of a second field within
// each group. The result is emitted as a single JSON object suitable for
// further processing or display.
package pivot

import (
	"encoding/json"
	"fmt"
)

// Pivotter groups log messages by groupKey and counts occurrences of each
// distinct value of valueKey within each group.
type Pivotter struct {
	groupKey string
	valueKey string
	// table[groupValue][valueValue] = count
	table map[string]map[string]int
}

// New creates a Pivotter that pivots on groupKey and counts valueKey.
// Both keys must be non-empty.
func New(groupKey, valueKey string) (*Pivotter, error) {
	if groupKey == "" {
		return nil, fmt.Errorf("pivot: groupKey must not be empty")
	}
	if valueKey == "" {
		return nil, fmt.Errorf("pivot: valueKey must not be empty")
	}
	return &Pivotter{
		groupKey: groupKey,
		valueKey: valueKey,
		table:    make(map[string]map[string]int),
	}, nil
}

// Add ingests a raw JSON log line. Lines that cannot be parsed or that are
// missing either key are silently skipped.
func (p *Pivotter) Add(line []byte) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(line, &m); err != nil {
		return
	}
	groupVal, ok := stringValue(m, p.groupKey)
	if !ok {
		return
	}
	valVal, ok := stringValue(m, p.valueKey)
	if !ok {
		return
	}
	if p.table[groupVal] == nil {
		p.table[groupVal] = make(map[string]int)
	}
	p.table[groupVal][valVal]++
}

// Result returns the pivot table as a JSON-encoded byte slice.
// The structure is: { "<groupValue>": { "<valueValue>": <count>, … }, … }
func (p *Pivotter) Result() ([]byte, error) {
	return json.Marshal(p.table)
}

// Reset clears all accumulated data.
func (p *Pivotter) Reset() {
	p.table = make(map[string]map[string]int)
}

// stringValue extracts a string value for key from a raw JSON map.
func stringValue(m map[string]json.RawMessage, key string) (string, bool) {
	raw, ok := m[key]
	if !ok {
		return "", false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	return s, true
}
