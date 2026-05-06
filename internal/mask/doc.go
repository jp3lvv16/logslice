// Package mask provides regex-based value masking for structured log fields.
//
// Each masking rule targets a specific JSON field and replaces portions of its
// string value that match a regular expression with a fixed mask string.
//
// Specs use the format:
//
//	field/pattern/mask
//
// For example:
//
//	email/[^@]+@/***@
//
// masks the local-part of an email address, turning "alice@example.com" into
// "***@example.com".
//
// Lines that are not valid JSON objects, or whose target field is not a string,
// are passed through unmodified.
package mask
