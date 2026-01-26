package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardDeleteCmd = &cobra.Command{
	Use:   "delete <card_number>",
	Short: "Delete a card",
	Long:  `Delete an existing card permanently`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.DeleteCard(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("deleting card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card #%d deleted successfully\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardDeleteCmd)
}
