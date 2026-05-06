// Package highlight provides lightweight ANSI terminal colouring for
// logslice's human-readable output modes.
//
// Usage:
//
//	p := highlight.New(isTTY)
//	fmt.Println(p.Level("error"), p.Key("msg"), p.Value("disk full"))
//
// When the Painter is created with enabled=false every method returns its
// input unchanged, making it safe to use unconditionally regardless of
// whether stdout is a terminal.
package highlight
