package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var pinListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pinned cards",
	Long:  `Retrieve and display the current user's pinned cards`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListPins(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListPins(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	cards, err := a.Client.GetMyPins(context.Background())
	if err != nil {
		return fmt.Errorf("fetching pinned cards: %w", err)
	}

	if len(cards) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No pinned cards")
		return nil
	}

	return ui.DisplayCards(cards)
}

func init() {
	pinCmd.AddCommand(pinListCmd)
}
