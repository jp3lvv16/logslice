// Package join provides a Joiner that correlates structured log lines from
// multiple sources by a shared key field.
//
// # Overview
//
// When processing logs from distributed systems it is common to have related
// events emitted separately — for example a request log and a response log
// that share a trace or request ID.  The Joiner buffers the first occurrence
// of a key and, when a second line carrying the same key arrives within the
// configured time window, merges the two objects into one and returns it.
//
// Lines whose key does not appear a second time within the window are evicted
// silently; call Flush to retrieve any remaining buffered entries at
// end-of-stream.
//
// # Example
//
//	j, _ := join.New("request_id", 5*time.Second)
//	for _, line := range lines {
//		out, _ := j.Add(line)
//		if out != nil {
//			fmt.Println(string(out))
//		}
//	}
//	for _, remaining := range j.Flush() { … }
package join
