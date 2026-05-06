// Package redact provides field-level value redaction for structured JSON log
// lines.
//
// # Usage
//
// Create a Redactor with the list of top-level field names that should be
// hidden and an optional replacement placeholder:
//
//	redactor, err := redact.New([]string{"password", "token"}, "[REDACTED]")
//	if err != nil { /* handle */ }
//
//	sanitised, err := redactor.Apply(rawJSONLine)
//
// Fields that are absent from the log line are silently skipped. Lines that
// cannot be parsed as JSON objects are passed through unchanged so that the
// pipeline is never interrupted by unexpected input.
package redact
