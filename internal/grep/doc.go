// Package grep provides full-text regular-expression search over raw JSON
// log lines.
//
// Unlike the field-level [filter] package, grep operates on the raw byte
// representation of each line before any JSON parsing takes place.  This
// makes it fast for simple keyword searches and useful as a first-pass
// filter in a pipeline.
//
// Basic usage:
//
//	g, err := grep.New(`"level":"error"`, false)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if g.Keep(rawLine) {
//		// process the line
//	}
//
// Inverted mode:
//
//	g, _ := grep.New(`debug`, true)  // keep lines that do NOT contain "debug"
package grep
