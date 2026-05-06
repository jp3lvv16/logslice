// Package tail provides functionality to follow a log file in real-time,
// emitting new lines as they are appended (similar to `tail -f`).
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Options controls the behaviour of Follow.
type Options struct {
	// PollInterval is how often to check for new data when the reader
	// reaches EOF. Defaults to 200ms if zero.
	PollInterval time.Duration
}

func (o Options) pollInterval() time.Duration {
	if o.PollInterval <= 0 {
		return 200 * time.Millisecond
	}
	return o.PollInterval
}

// Follow opens the file at path and sends each new line to lines as it
// arrives. It blocks until ctx is cancelled or a read error occurs.
// The returned error is nil when ctx is cancelled normally.
func Follow(ctx context.Context, path string, lines chan<- string, opts Options) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end so we only emit lines written after Follow is called.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(opts.pollInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					// Strip trailing newline before sending.
					if len(line) > 0 && line[len(line)-1] == '\n' {
						line = line[:len(line)-1]
					}
					select {
					case lines <- line:
					case <-ctx.Done():
						return nil
					}
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}
		}
	}
}
