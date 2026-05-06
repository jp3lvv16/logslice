// Package transform applies field renaming and value rewriting to structured
// JSON log records before they are written to output.
//
// Two spec formats are supported:
//
//	"oldKey=newKey"       — rename a field, preserving its value.
//	"key:oldVal=newVal"   — rewrite a field's string value when it matches
//	                        oldVal, leaving the key name unchanged.
//
// Specs are evaluated in the order they are supplied. Lines that cannot be
// parsed as JSON objects are passed through unmodified.
package transform
