// Package throttle implements a sliding-window line throttler for log
// pipelines.
//
// A Throttler is constructed with a maximum line count and a window
// duration. Within each window period the throttler allows up to that
// many lines through; any additional lines are silently dropped and
// counted so callers can surface the loss in summary statistics.
//
// Example
//
//	th, err := throttle.New(100, time.Second)
//	if err != nil { ... }
//	if th.Allow(msg) {
//	    // forward msg downstream
//	}
//	fmt.Printf("dropped %d lines\n", th.Dropped())
package throttle
