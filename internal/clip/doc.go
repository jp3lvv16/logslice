// Package clip implements a sliding-window byte-budget gate for structured
// log streams.
//
// # Overview
//
// A Clipper is initialised with a maximum byte allowance and a rolling time
// window duration.  Each call to Allow encodes the message size, evicts
// entries that have aged out of the window, and then either records the
// message and returns true (within budget) or returns false (budget
// exhausted).
//
// # Usage
//
//	c, err := clip.New(1_000_000, time.Minute) // 1 MB per minute
//	if err != nil { … }
//
//	if c.Allow(rawMsg) {
//		// forward the message
//	}
//
// Setting maxBytes to 0 or a negative value disables clipping entirely —
// every message is allowed through unconditionally.
package clip
