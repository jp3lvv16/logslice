package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/user/logslice/internal/throttle"
)

// registerThrottleFlags wires --throttle-max and --throttle-window into fs.
func registerThrottleFlags(fs *flag.FlagSet, max *int, window *time.Duration) {
	fs.IntVar(max, "throttle-max", 0,
		"maximum log lines per throttle window (0 = unlimited)")
	fs.DurationVar(window, "throttle-window", time.Second,
		"duration of the throttle window (e.g. 1s, 500ms)")
}

// buildThrottler returns a configured Throttler or nil when throttling is
// disabled (max <= 0).
func buildThrottler(max int, window time.Duration) (*throttle.Throttler, error) {
	if max <= 0 {
		return nil, nil
	}
	th, err := throttle.New(max, window)
	if err != nil {
		return nil, fmt.Errorf("throttle: %w", err)
	}
	return th, nil
}
