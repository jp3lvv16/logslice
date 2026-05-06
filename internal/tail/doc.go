// Package tail implements real-time file following for logslice.
//
// It exposes a single function, Follow, which opens a file, seeks to its
// current end, and continuously emits newly appended lines on a channel.
// This enables logslice to operate in a streaming "live" mode analogous to
// `tail -f`, while still passing each line through the normal pipeline of
// time-range filtering, field matching, deduplication, and output formatting.
//
// Usage:
//
//	lines := make(chan string, 64)
//	go tail.Follow(ctx, "/var/log/app.log", lines, tail.Options{})
//	for line := range lines {
//		// process line
//	}
package tail
