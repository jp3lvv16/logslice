package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

func writeLine(t *testing.T, f *os.File, s string) {
	t.Helper()
	_, err := f.WriteString(s + "\n")
	if err != nil {
		t.Fatalf("writeLine: %v", err)
	}
}

func TestFollow_EmitsNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	lines := make(chan string, 10)
	opts := tail.Options{PollInterval: 20 * time.Millisecond}

	go func() {
		if err := tail.Follow(ctx, f.Name(), lines, opts); err != nil {
			t.Errorf("Follow: %v", err)
		}
	}()

	// Give Follow time to open and seek to end.
	time.Sleep(50 * time.Millisecond)

	wantLines := []string{`{"msg":"hello"}`, `{"msg":"world"}`}
	for _, l := range wantLines {
		writeLine(t, f, l)
	}

	for _, want := range wantLines {
		select {
		case got := <-lines:
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		case <-ctx.Done():
			t.Fatalf("timed out waiting for line %q", want)
		}
	}
}

func TestFollow_CancelStops(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-cancel-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	lines := make(chan string, 1)
	opts := tail.Options{PollInterval: 20 * time.Millisecond}

	done := make(chan error, 1)
	go func() {
		done <- tail.Follow(ctx, f.Name(), lines, opts)
	}()

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("expected nil error after cancel, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Follow did not stop after context cancel")
	}
}

func TestFollow_MissingFile(t *testing.T) {
	ctx := context.Background()
	lines := make(chan string, 1)
	err := tail.Follow(ctx, "/nonexistent/path/file.log", lines, tail.Options{})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
