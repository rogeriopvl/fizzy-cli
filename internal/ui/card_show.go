package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/rogeriopvl/fizzy/internal/api"
)

func DisplayCard(card *api.Card) error {
	boldStyle := lipgloss.NewStyle().Bold(true)
	dimStyle := lipgloss.NewStyle().Bold(true).Faint(true)
	fmt.Printf("%s\n", boldStyle.Render(fmt.Sprintf("%s (#%d)", card.Title, card.Number)))
	fmt.Println("─────────────────────────────────────")
	fmt.Printf("%s %s\n", dimStyle.Render("Description:"), card.Description)
	fmt.Printf("%s %v\n", dimStyle.Render("Tags:"), card.Tags)
	fmt.Printf("%s %v\n", dimStyle.Render("Golden:"), card.Golden)
	fmt.Printf("%s %v\n", dimStyle.Render("Closed:"), card.Closed)
	if card.Column != nil {
		fmt.Printf("%s %s\n", dimStyle.Render("Column:"), card.Column.Name)
	}
	fmt.Printf("%s %s\n", dimStyle.Render("Status:"), card.Status)
	fmt.Printf("%s %s\n", dimStyle.Render("Created:"), card.CreatedAt)
	fmt.Printf("%s %s\n", dimStyle.Render("Last Active:"), card.LastActiveAt)
	fmt.Printf("%s %s\n", dimStyle.Render("URL:"), card.URL)
	if len(card.Steps) > 0 {
		fmt.Printf("%s\n", dimStyle.Render("Steps:"))
		for _, step := range card.Steps {
			checkmark := "☐"
			if step.Completed {
				checkmark = "☑"
			}
			fmt.Printf("  %s %s\n", checkmark, step.Content)
		}
	}
	return nil
}
