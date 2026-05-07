// Package coalesce provides a field-coalescing processor for structured log
// records.
//
// Different log producers often use inconsistent field names for the same
// semantic value — e.g. "msg", "message", "text", or "body" for the human-
// readable log line. Coalescer normalises these into a single canonical field
// by scanning a prioritised list of candidate fields and promoting the first
// non-empty value.
//
// Example
//
//	c, err := coalesce.New("message", []string{"msg", "text", "body"}, false)
//	if err != nil { ... }
//	out, err := c.Apply(raw)
package coalesce
