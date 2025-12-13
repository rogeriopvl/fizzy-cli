package ui

import "github.com/charmbracelet/lipgloss"

// DisplayID returns a formatted ID string with greyed out styling
func DisplayID(id string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Grey color
		Render("id: " + id)
}
