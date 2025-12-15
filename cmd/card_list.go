package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var cardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cards",
	Long:  `Retrieve and display all cards from Fizzy`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListCards(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListCards(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if a.Config.SelectedBoard == "" {
		return fmt.Errorf("no board selected")
	}

	filters := api.CardFilters{
		BoardIDs: []string{a.Config.SelectedBoard},
	}

	cards, err := a.Client.GetCards(context.Background(), filters)
	if err != nil {
		return fmt.Errorf("fetching cards: %w", err)
	}

	if len(cards) == 0 {
		fmt.Println("No cards found")
		return nil
	}

	return ui.DisplayCards(cards)
}

func init() {
	cardCmd.AddCommand(cardListCmd)
}
