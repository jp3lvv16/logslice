package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/user/logslice/internal/join"
)

// joinConfig holds the parsed join flags.
type joinConfig struct {
	key    string
	window time.Duration
}

// registerJoinFlags registers --join-key and --join-window onto fs.
func registerJoinFlags(fs *flag.FlagSet, cfg *joinConfig) {
	fs.StringVar(&cfg.key, "join-key", "", "field name to join log lines on (empty = disabled)")
	fs.DurationVar(&cfg.window, "join-window", 5*time.Second, "maximum time between correlated lines")
}

// buildJoiner returns a configured *join.Joiner when joining is enabled, or
// nil when the key flag is empty. An error is returned for invalid configs.
func buildJoiner(cfg joinConfig) (*join.Joiner, error) {
	if cfg.key == "" {
		return nil, nil
	}
	j, err := join.New(cfg.key, cfg.window)
	if err != nil {
		return nil, fmt.Errorf("join: %w", err)
	}
	return j, nil
}
