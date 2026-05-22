package formatter

import (
	"fmt"
	"time"
)

// Helper functions shared across formatters

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// humanizeDuration renders a duration as a short relative string.
// Examples: "just now", "5m ago", "3h ago", "2d ago".
func humanizeDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dw ago", int(d.Hours()/(24*7)))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(d.Hours()/(24*365)))
	}
}

// relativeTime renders how long ago t occurred.
// It returns an empty string for the zero time.
func relativeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return humanizeDuration(time.Since(t))
}

// formatTimestamp renders an absolute timestamp with a relative suffix.
// Example: "2025-08-30 10:00:00 (3h ago)".
func formatTimestamp(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	abs := t.Format("2006-01-02 15:04:05")
	if rel := relativeTime(t); rel != "" {
		return abs + " (" + rel + ")"
	}
	return abs
}

// dashIfEmpty returns "-" when s is empty, keeping table columns aligned.
func dashIfEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
