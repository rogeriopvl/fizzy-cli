package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var cardReactionListCmd = &cobra.Command{
	Use:   "list <card_number>",
	Short: "List reactions on a card",
	Long:  `Retrieve and display all reactions (boosts) on a card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListCardReactions(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListCardReactions(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	reactions, err := a.Client.GetCardReactions(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("fetching reactions: %w", err)
	}

	if len(reactions) == 0 {
		fmt.Println("No reactions found")
		return nil
	}

	return ui.DisplayReactions(reactions)
}

func init() {
	cardReactionCmd.AddCommand(cardReactionListCmd)
}
