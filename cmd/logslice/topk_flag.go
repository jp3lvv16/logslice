package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/logslice/internal/topk"
)

// registerTopKFlags adds --topk-field and --topk-k to the provided FlagSet.
func registerTopKFlags(fs *flag.FlagSet, field *string, k *int) {
	fs.StringVar(field, "topk-field", "", "JSON field to track for top-K frequency report")
	fs.IntVar(k, "topk-k", 10, "number of top values to display in the frequency report")
}

// buildTopKTracker constructs a *topk.Tracker when field is non-empty.
// Returns nil, nil when the feature is disabled.
func buildTopKTracker(field string, k int) (*topk.Tracker, error) {
	if field == "" {
		return nil, nil
	}
	tr, err := topk.New(field, k)
	if err != nil {
		return nil, fmt.Errorf("--topk-field: %w", err)
	}
	return tr, nil
}

// printTopK writes the frequency table to stderr when tracker is non-nil.
func printTopK(tr *topk.Tracker, field string) {
	if tr == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "\ntop-%d values for %q:\n", cap(tr.Top())+len(tr.Top()), field)
	for i, e := range tr.Top() {
		fmt.Fprintf(os.Stderr, "  %2d. %-30s %d\n", i+1, e.Value, e.Count)
	}
}
