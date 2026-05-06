// Package context provides a pipeline context that carries per-line
// metadata (line number, source file, ingestion timestamp) through the
// logslice processing stages.
package context

import (
	"encoding/json"
	"time"
)

// Meta holds metadata injected alongside every log line as it travels
// through the processing pipeline.
type Meta struct {
	Line      int64     `json:"_line"`
	Source    string    `json:"_source,omitempty"`
	IngestedAt time.Time `json:"_ingested_at"`
}

// Enrich injects the Meta fields into a raw JSON message, returning the
// augmented bytes. Fields that already exist in msg are NOT overwritten.
func Enrich(msg json.RawMessage, m Meta) (json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return msg, err
	}

	if _, ok := obj["_line"]; !ok {
		v, _ := json.Marshal(m.Line)
		obj["_line"] = v
	}
	if m.Source != "" {
		if _, ok := obj["_source"]; !ok {
			v, _ := json.Marshal(m.Source)
			obj["_source"] = v
		}
	}
	if _, ok := obj["_ingested_at"]; !ok {
		v, _ := json.Marshal(m.IngestedAt.UTC().Format(time.RFC3339Nano))
		obj["_ingested_at"] = v
	}

	return json.Marshal(obj)
}

// Strip removes all Meta fields (_line, _source, _ingested_at) from a
// JSON message before final output.
func Strip(msg json.RawMessage) (json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(msg, &obj); err != nil {
		return msg, err
	}
	delete(obj, "_line")
	delete(obj, "_source")
	delete(obj, "_ingested_at")
	return json.Marshal(obj)
}
