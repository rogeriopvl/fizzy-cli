package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var cardShowCmd = &cobra.Command{
	Use:   "show <card_id>",
	Short: "Show card details",
	Long:  `Retrieve and display details for a specific card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowCard(cmd *cobra.Command, cardID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	card, err := a.Client.GetCard(context.Background(), cardID)
	if err != nil {
		return fmt.Errorf("fetching card: %w", err)
	}

	return ui.DisplayCard(card)
}

func init() {
	cardCmd.AddCommand(cardShowCmd)
}
