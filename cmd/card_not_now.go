package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardNotNowCmd = &cobra.Command{
	Use:   "not-now <card_number>",
	Short: "Move a card to Not Now status",
	Long:  `Move an existing card to the "Not Now" status`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleNotNowCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleNotNowCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.PostCardNotNow(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("moving card to not now: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card #%d moved to Not Now successfully\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardNotNowCmd)
}
