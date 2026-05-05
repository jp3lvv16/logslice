// Package filter provides field-pattern filtering for structured log entries.
//
// A Filter is created from one or more spec strings of the form:
//
//	"field=pattern"
//
// where pattern is a Go regular expression. A log entry (represented as a
// map[string]string) matches the filter only when every rule's field exists
// in the entry and the field's value matches the corresponding pattern.
//
// Example usage:
//
//	f, err := filter.New([]string{"level=error", "service=auth.*"})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	entry := map[string]string{
//		"level":   "error",
//		"service": "auth-api",
//		"msg":     "token expired",
//	}
//
//	if f.Match(entry) {
//		fmt.Println("entry matches filter")
//	}
package filter
