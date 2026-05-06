package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestColorize_Disabled(t *testing.T) {
	p := highlight.New(false)
	got := p.Colorize(highlight.Red, "hello")
	if got != "hello" {
		t.Fatalf("expected plain text, got %q", got)
	}
}

func TestColorize_Enabled(t *testing.T) {
	p := highlight.New(true)
	got := p.Colorize(highlight.Red, "hello")
	if !strings.Contains(got, "hello") {
		t.Fatal("output must contain original string")
	}
	if !strings.HasPrefix(got, "\033[") {
		t.Fatal("output must start with ANSI escape")
	}
	if !strings.HasSuffix(got, highlight.Reset) {
		t.Fatal("output must end with reset code")
	}
}

func TestColorize_EmptyString(t *testing.T) {
	p := highlight.New(true)
	if got := p.Colorize(highlight.Green, ""); got != "" {
		t.Fatalf("empty string should pass through unchanged, got %q", got)
	}
}

func TestLevel_Colours(t *testing.T) {
	p := highlight.New(true)
	cases := []struct {
		input    string
		wantCode string
	}{
		{"error", highlight.Red},
		{"FATAL", highlight.Red},
		{"warn", highlight.Yellow},
		{"WARNING", highlight.Yellow},
		{"info", highlight.Green},
		{"INFO", highlight.Green},
		{"debug", highlight.Blue},
		{"TRACE", highlight.Blue},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			out := p.Level(tc.input)
			if !strings.Contains(out, tc.wantCode) {
				t.Errorf("Level(%q) = %q, want ANSI code %q", tc.input, out, tc.wantCode)
			}
		})
	}
}

func TestLevel_Unknown(t *testing.T) {
	p := highlight.New(true)
	got := p.Level("notice")
	// unknown levels should pass through without ANSI codes
	if strings.Contains(got, "\033[") {
		t.Fatalf("unknown level should not be coloured, got %q", got)
	}
}

func TestKey_Value_Bold(t *testing.T) {
	p := highlight.New(true)
	if !strings.Contains(p.Key("ts"), highlight.Cyan) {
		t.Error("Key should use Cyan")
	}
	if !strings.Contains(p.Value("42"), highlight.Yellow) {
		t.Error("Value should use Yellow")
	}
	if !strings.Contains(p.Bold("hi"), highlight.Bold) {
		t.Error("Bold should use Bold code")
	}
}
