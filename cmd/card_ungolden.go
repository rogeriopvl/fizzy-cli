package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardUngoldenCmd = &cobra.Command{
	Use:   "ungolden <card_number>",
	Short: "Remove golden status from a card",
	Long:  `Remove golden status from an existing card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUngoldenCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUngoldenCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err = a.Client.UnmarkCardGolden(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("removing golden status: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card #%d golden status removed\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardUngoldenCmd)
}
