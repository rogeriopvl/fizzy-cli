package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardCloseCmd = &cobra.Command{
	Use:   "close <card_number>",
	Short: "Close a card",
	Long:  `Close an existing card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCloseCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCloseCard(cmd *cobra.Command, cardNumber string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	_, err = a.Client.PostCardsClosure(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("closing card: %w", err)
	}

	fmt.Printf("✓ Card #%d closed successfully\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardCloseCmd)
}
