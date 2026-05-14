// Package routing implements field-value-based routing for structured log lines.
//
// A Router is constructed from a list of routing specs, each of the form:
//
//	"field=pattern:dest"
//
// Where:
//   - field   is the JSON key to inspect
//   - pattern is a regular expression matched against the field's string value
//   - dest    is the symbolic destination name returned on a match
//
// Rules are evaluated in order; the first match wins. If no rule matches the
// configured defaultDest is returned. Non-JSON lines and lines missing the
// target field are also routed to defaultDest.
//
// Example:
//
//	router, _ := routing.New(
//	    []string{"level=error:errors", "level=warn:warnings"},
//	    "default",
//	)
//	dest := router.Route(line) // "errors", "warnings", or "default"
package routing
