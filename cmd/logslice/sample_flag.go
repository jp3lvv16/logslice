package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/user/logslice/internal/sample"
)

// registerSampleFlags adds --sample-size to fs.
func registerSampleFlags(fs *flag.FlagSet, size *int) {
	fs.IntVar(size, "sample-size", 0,
		"reservoir sample: keep at most N randomly selected lines (0 = disabled)")
}

// buildSampler constructs a *sample.Sampler when size > 0, otherwise returns nil.
func buildSampler(size int) (*sample.Sampler, error) {
	if size <= 0 {
		return nil, nil
	}
	s, err := sample.New(size, time.Now().UnixNano())
	if err != nil {
		return nil, fmt.Errorf("--sample-size: %w", err)
	}
	return s, nil
}
