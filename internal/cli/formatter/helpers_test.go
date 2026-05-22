package formatter

import (
	"strings"
	"testing"
	"time"
)

func TestHumanizeDuration(t *testing.T) {
	cases := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"seconds", 30 * time.Second, "just now"},
		{"minutes", 5 * time.Minute, "5m ago"},
		{"hours", 3 * time.Hour, "3h ago"},
		{"days", 2 * 24 * time.Hour, "2d ago"},
		{"weeks", 3 * 7 * 24 * time.Hour, "3w ago"},
		{"months", 60 * 24 * time.Hour, "2mo ago"},
		{"years", 400 * 24 * time.Hour, "1y ago"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := humanizeDuration(tc.duration)
			if got != tc.expected {
				t.Errorf("humanizeDuration(%v) = %q, want %q", tc.duration, got, tc.expected)
			}
		})
	}
}

func TestRelativeTimeZeroValue(t *testing.T) {
	if got := relativeTime(time.Time{}); got != "" {
		t.Errorf("relativeTime(zero) = %q, want empty string", got)
	}
}

func TestRelativeTimeRecent(t *testing.T) {
	got := relativeTime(time.Now().Add(-2 * time.Hour))
	if got != "2h ago" {
		t.Errorf("relativeTime(2h back) = %q, want %q", got, "2h ago")
	}
}

func TestFormatTimestamp(t *testing.T) {
	recent := time.Now().Add(-3 * time.Hour)
	got := formatTimestamp(recent)
	if !strings.Contains(got, "(3h ago)") {
		t.Errorf("formatTimestamp = %q, want it to contain relative suffix", got)
	}
	if !strings.Contains(got, recent.Format("2006-01-02")) {
		t.Errorf("formatTimestamp = %q, want it to contain absolute date", got)
	}

	if got := formatTimestamp(time.Time{}); got != "" {
		t.Errorf("formatTimestamp(zero) = %q, want empty string", got)
	}
}

func TestDashIfEmpty(t *testing.T) {
	if got := dashIfEmpty(""); got != "-" {
		t.Errorf("dashIfEmpty(\"\") = %q, want %q", got, "-")
	}
	if got := dashIfEmpty("high"); got != "high" {
		t.Errorf("dashIfEmpty(\"high\") = %q, want %q", got, "high")
	}
}
