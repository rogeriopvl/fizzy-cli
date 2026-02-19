package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardGoldenCmd = &cobra.Command{
	Use:   "golden <card_number>",
	Short: "Mark a card as golden",
	Long:  `Mark an existing card as golden`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleGoldenCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleGoldenCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err = a.Client.MarkCardGolden(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("marking card as golden: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card #%d marked as golden\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardGoldenCmd)
}
