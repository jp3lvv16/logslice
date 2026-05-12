// Package window implements a sliding time-window buffer for structured log
// lines.
//
// Lines are grouped into fixed-duration buckets determined by each message's
// parsed timestamp.  When a new message falls outside the current bucket the
// completed bucket is flushed and a fresh one is started.
//
// Typical usage:
//
//	w, _ := window.New(5 * time.Minute)
//	for _, entry := range entries {
//		if group, flushed := w.Add(entry.Time, entry.Raw); flushed {
//			process(group)
//		}
//	}
//	if tail := w.Flush(); tail != nil {
//		process(tail)
//	}
package window
