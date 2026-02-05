package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardReactionDeleteCmd = &cobra.Command{
	Use:   "delete <card_number> <reaction_id>",
	Short: "Delete a reaction from a card",
	Long:  `Remove your reaction (boost) from a card`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteCardReaction(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteCardReaction(cmd *cobra.Command, cardNumber, reactionID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.DeleteCardReaction(context.Background(), cardNum, reactionID)
	if err != nil {
		return fmt.Errorf("deleting reaction: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Reaction deleted successfully\n")
	return nil
}

func init() {
	cardReactionCmd.AddCommand(cardReactionDeleteCmd)
}
