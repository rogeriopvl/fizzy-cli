package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/rogeriopvl/fizzy-cli/internal/api"
	"github.com/rogeriopvl/fizzy-cli/internal/colors"
)

func DisplayColumns(columns []api.Column) error {
	for _, column := range columns {
		colorName := column.Color.Name
		colorDef := colors.ByName(colorName)

		termColor := lipgloss.Color("7") // default to white
		if colorDef != nil {
			termColor = colorDef.TermColor
		}

		styledName := lipgloss.NewStyle().
			Foreground(termColor).
			Render(column.Name)

		fmt.Printf("%s (%s)\n", styledName, DisplayID(column.ID))
	}
	return nil
}
