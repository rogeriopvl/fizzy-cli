package ui

import "github.com/charmbracelet/lipgloss"

// DisplayMeta returns a formatted metadata string with greyed out styling
func DisplayMeta(label string, value string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Grey color
		Render(label + ": " + value)
}

func DisplayID(id string) string {
	return DisplayMeta("id", id)
}
