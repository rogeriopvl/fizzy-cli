package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardReopenCmd = &cobra.Command{
	Use:   "reopen <card_number>",
	Short: "Reopen a card",
	Long:  `Reopen an existing closed card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleReopenCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleReopenCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.DeleteCardsClosure(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("reopening card: %w", err)
	}

	fmt.Printf("✓ Card #%d reopened successfully\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardReopenCmd)
}
