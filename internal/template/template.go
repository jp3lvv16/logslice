// Package template provides a simple field-interpolation formatter that
// renders a Go text/template string against each JSON log message.
//
// Example spec: "{{.level}} {{.msg}} ({{.service}})"
package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

// Renderer applies a Go template to each raw JSON log line.
type Renderer struct {
	tmpl *template.Template
}

// New compiles the template string and returns a Renderer.
// Returns an error if the template cannot be parsed.
func New(tmplStr string) (*Renderer, error) {
	if tmplStr == "" {
		return nil, fmt.Errorf("template: spec must not be empty")
	}
	t, err := template.New("logslice").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template: %w", err)
	}
	return &Renderer{tmpl: t}, nil
}

// Apply renders the template against the decoded JSON message and returns the
// resulting string. Fields missing from the message are rendered as an empty
// string (missingkey=zero). Returns an error if msg is not valid JSON or if
// template execution fails.
func (r *Renderer) Apply(msg json.RawMessage) (string, error) {
	var data map[string]any
	if err := json.Unmarshal(msg, &data); err != nil {
		return "", fmt.Errorf("template: unmarshal: %w", err)
	}

	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template: execute: %w", err)
	}
	return buf.String(), nil
}
