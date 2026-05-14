package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/user/logslice/internal/routing"
)

// routingFlag is a multi-value flag that accumulates routing specs.
type routingFlag []string

func (f *routingFlag) String() string { return strings.Join(*f, ",") }
func (f *routingFlag) Set(v string) error {
	*f = append(*f, strings.TrimSpace(v))
	return nil
}

// registerRoutingFlags wires the --route and --route-default flags into fs.
func registerRoutingFlags(fs *flag.FlagSet, specs *routingFlag, defaultDest *string) {
	fs.Var(specs, "route",
		`route lines by field pattern: "field=regex:dest" (repeatable)`)
	fs.StringVar(defaultDest, "route-default", "default",
		"destination name used when no routing rule matches")
}

// buildRouter constructs a *routing.Router from the accumulated specs.
// Returns nil and no error when no specs are provided.
func buildRouter(specs []string, defaultDest string) (*routing.Router, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	r, err := routing.New(specs, defaultDest)
	if err != nil {
		return nil, fmt.Errorf("--route: %w", err)
	}
	return r, nil
}
