// Package truncate caps the length of string field values inside structured
// JSON log entries.
//
// # Usage
//
//	tr, err := truncate.New(120, "…")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	truncated, err := tr.Apply(rawMsg)
//
// Any string field whose byte length exceeds maxLen is replaced with the first
// maxLen bytes followed by the configured suffix. All other field types (numbers,
// booleans, nested objects, arrays) are left untouched.
//
// Passing maxLen = 0 disables truncation and returns the original message as-is.
package truncate
