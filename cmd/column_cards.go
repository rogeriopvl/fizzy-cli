package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var columnCardsCmd = &cobra.Command{
	Use:   "cards <column_id>",
	Short: "List cards in a column",
	Long:  `Retrieve and display cards in a specific column of the selected board`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListColumnCards(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListColumnCards(cmd *cobra.Command, columnID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	opts := &fizzy.ListOptions{}
	if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
		opts.Limit = limit
	}

	cards, err := a.Client.GetColumnCards(context.Background(), columnID, opts)
	if err != nil {
		return fmt.Errorf("fetching column cards: %w", err)
	}

	if len(cards) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No cards found")
		return nil
	}

	return ui.DisplayCards(cards)
}

func init() {
	columnCardsCmd.Flags().IntP("limit", "l", 0, "Maximum number of cards to return (0 = no limit)")
	columnCmd.AddCommand(columnCardsCmd)
}
