package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardAssignCmd = &cobra.Command{
	Use:   "assign <card_number> <user_id>",
	Short: "Assign a user to a card",
	Long: `Assign or unassign a user to/from a card.

Use "me" as the user_id to assign the card to yourself.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleAssignCard(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleAssignCard(cmd *cobra.Command, cardNumber, userID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if userID == "me" {
		if a.Config.CurrentUserID == "" {
			return fmt.Errorf("current user ID not available, please run 'fizzy login' first")
		}
		userID = a.Config.CurrentUserID
	}

	_, err = a.Client.PostCardAssignments(context.Background(), cardNum, userID)
	if err != nil {
		return fmt.Errorf("assigning card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card #%d assignment toggled for user %s\n", cardNum, userID)
	return nil
}

func init() {
	cardCmd.AddCommand(cardAssignCmd)
}
