package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var cardUnpinCmd = &cobra.Command{
	Use:   "unpin <card_number>",
	Short: "Unpin a card",
	Long:  `Remove a card from your pinned cards list`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUnpinCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUnpinCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if err := a.Client.UnpinCard(context.Background(), cardNum); err != nil {
		return fmt.Errorf("unpinning card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card #%d unpinned\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardUnpinCmd)
}
