package ui

import (
	"github.com/charmbracelet/lipgloss"
	"time"
)

// DisplayMeta returns a formatted metadata string with greyed out styling
func DisplayMeta(label string, value string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Grey color
		Render(label + ": " + value)
}

func DisplayID(id string) string {
	return DisplayMeta("id", id)
}

// FormatTime converts an RFC3339 timestamp string to a human-readable format.
// If parsing fails, returns the original string.
func FormatTime(timeStr string) string {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	return t.Format("2006-01-02 15:04:05")
}
