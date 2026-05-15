// Package normalize provides field-level string normalization for structured
// JSON log messages.
//
// # Overview
//
// A Normalizer is constructed with a list of specs, each of the form:
//
//	"field:mode"
//
// where mode is one of:
//
//	lower  – convert the field value to lower-case
//	upper  – convert the field value to upper-case
//	title  – convert the field value to title-case
//	trim   – strip leading and trailing whitespace
//
// # Example
//
//	n, err := normalize.New([]string{"level:lower", "msg:trim"})
//	if err != nil { ... }
//	out, err := n.Apply(msg)
//
// Fields that are absent or whose value is not a JSON string are silently
// skipped so that normalization never corrupts numeric or boolean fields.
package normalize
