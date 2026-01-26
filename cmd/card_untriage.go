package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardUntriagedCmd = &cobra.Command{
	Use:   "untriage <card_number>",
	Short: "Send a card back to triage",
	Long:  `Send an existing card back to the triage column`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUntriagedCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUntriagedCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.DeleteCardTriage(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("sending card back to triage: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card #%d sent back to triage successfully\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardUntriagedCmd)
}
