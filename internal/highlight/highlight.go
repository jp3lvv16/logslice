// Package highlight provides ANSI colour helpers for terminal output.
package highlight

import (
	"fmt"
	"strings"
)

// Colour codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Painter wraps output with ANSI codes when colour is enabled.
type Painter struct {
	enabled bool
}

// New returns a Painter. When enabled is false all methods return plain text.
func New(enabled bool) *Painter {
	return &Painter{enabled: enabled}
}

// Colorize wraps s with the given ANSI code when colour is enabled.
func (p *Painter) Colorize(code, s string) string {
	if !p.enabled || s == "" {
		return s
	}
	return fmt.Sprintf("%s%s%s", code, s, Reset)
}

// Key returns a field key styled for display.
func (p *Painter) Key(k string) string {
	return p.Colorize(Cyan, k)
}

// Value returns a field value styled for display.
func (p *Painter) Value(v string) string {
	return p.Colorize(Yellow, v)
}

// Level returns the log level string coloured by severity.
func (p *Painter) Level(level string) string {
	if !p.enabled {
		return level
	}
	switch strings.ToLower(level) {
	case "error", "err", "fatal", "crit":
		return p.Colorize(Red, level)
	case "warn", "warning":
		return p.Colorize(Yellow, level)
	case "info":
		return p.Colorize(Green, level)
	case "debug", "trace":
		return p.Colorize(Blue, level)
	default:
		return level
	}
}

// Bold returns s in bold when colour is enabled.
func (p *Painter) Bold(s string) string {
	return p.Colorize(Bold, s)
}
