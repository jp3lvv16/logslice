// Package output handles formatting and writing of matched log entries.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for log entries.
type Format int

const (
	// FormatJSON outputs each entry as compact JSON.
	FormatJSON Format = iota
	// FormatPretty outputs each entry as indented JSON.
	FormatPretty
	// FormatText outputs selected fields as key=value pairs.
	FormatText
)

// Writer wraps an io.Writer and writes log entries in the configured format.
type Writer struct {
	out    io.Writer
	format Format
	fields []string
}

// New creates a new Writer with the given output destination, format, and
// optional field projection. If fields is empty all fields are written.
func New(out io.Writer, format Format, fields []string) *Writer {
	return &Writer{out: out, format: format, fields: fields}
}

// Write formats a single log entry (represented as a raw JSON message) and
// writes it to the underlying writer. It returns any write error.
func (w *Writer) Write(raw json.RawMessage) error {
	var entry map[string]interface{}
	if err := json.Unmarshal(raw, &entry); err != nil {
		return fmt.Errorf("output: unmarshal: %w", err)
	}

	if len(w.fields) > 0 {
		entry = project(entry, w.fields)
	}

	switch w.format {
	case FormatPretty:
		b, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			return fmt.Errorf("output: marshal: %w", err)
		}
		_, err = fmt.Fprintln(w.out, string(b))
		return err
	case FormatText:
		parts := make([]string, 0, len(entry))
		for k, v := range entry {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
		_, err := fmt.Fprintln(w.out, strings.Join(parts, " "))
		return err
	default: // FormatJSON
		b, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("output: marshal: %w", err)
		}
		_, err = fmt.Fprintln(w.out, string(b))
		return err
	}
}

// project returns a new map containing only the specified keys.
func project(entry map[string]interface{}, fields []string) map[string]interface{} {
	out := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		if v, ok := entry[f]; ok {
			out[f] = v
		}
	}
	return out
}
