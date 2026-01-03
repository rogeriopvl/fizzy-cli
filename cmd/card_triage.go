package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardTriageCmd = &cobra.Command{
	Use:   "triage <card_number> <column_id>",
	Short: "Move a card from triage into a column",
	Long:  `Move a card from triage into a specified column`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleTriageCard(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleTriageCard(cmd *cobra.Command, cardNumber string, columnID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.PostCardTriage(context.Background(), cardNum, columnID)
	if err != nil {
		return fmt.Errorf("triaging card: %w", err)
	}

	fmt.Printf("✓ Card #%d moved to column successfully\n", cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardTriageCmd)
}
