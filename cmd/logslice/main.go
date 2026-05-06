package main

import (
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/reader"
	"github.com/user/logslice/internal/timerange"
)

func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	cfg, err := parseFlags(args, stderr)
	if err != nil {
		return 2
	}

	tr, err := timerange.Parse(cfg.from, cfg.to)
	if err != nil {
		fmt.Fprintf(stderr, "error: invalid time range: %v\n", err)
		return 1
	}

	f, err := filter.New(cfg.filters)
	if err != nil {
		fmt.Fprintf(stderr, "error: invalid filter: %v\n", err)
		return 1
	}

	w, err := output.New(stdout, cfg.format, cfg.fields)
	if err != nil {
		fmt.Fprintf(stderr, "error: invalid output config: %v\n", err)
		return 1
	}

	src := stdin
	if cfg.file != "" {
		f2, err := os.Open(cfg.file)
		if err != nil {
			fmt.Fprintf(stderr, "error: cannot open file: %v\n", err)
			return 1
		}
		defer f2.Close()
		src = f2
	}

	if err := reader.Read(src, tr, f, w); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
