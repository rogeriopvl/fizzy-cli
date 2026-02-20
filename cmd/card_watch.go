package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardWatchCmd = &cobra.Command{
	Use:   "watch <card_number>",
	Short: "Subscribe to card notifications",
	Long:  `Subscribe to notifications for an existing card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleWatchCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleWatchCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err = a.Client.WatchCard(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("watching card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Now watching card #%d\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardWatchCmd)
}
