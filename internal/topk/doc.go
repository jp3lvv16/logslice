// Package topk provides a frequency tracker that records how often each
// distinct value of a chosen JSON field appears across a stream of log lines
// and returns the top-K most frequent values on demand.
//
// Typical usage:
//
//	tr, err := topk.New("level", 5)
//	if err != nil { ... }
//	for _, line := range lines {
//	    tr.Add(line)
//	}
//	for _, e := range tr.Top() {
//	    fmt.Printf("%s: %d\n", e.Value, e.Count)
//	}
package topk
