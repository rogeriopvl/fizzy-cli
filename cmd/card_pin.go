package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var cardPinCmd = &cobra.Command{
	Use:   "pin <card_number>",
	Short: "Pin a card",
	Long:  `Pin a card so it appears in your pinned cards list`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handlePinCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handlePinCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if err := a.Client.PinCard(context.Background(), cardNum); err != nil {
		return fmt.Errorf("pinning card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card #%d pinned\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardPinCmd)
}
