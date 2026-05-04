package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var cardImageDeleteCmd = &cobra.Command{
	Use:   "delete <card_number>",
	Short: "Delete a card's image",
	Long:  `Remove the image attached to a card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteCardImage(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteCardImage(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if err := a.Client.DeleteCardImage(context.Background(), cardNum); err != nil {
		return fmt.Errorf("deleting card image: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Image removed from card #%d\n", cardNum)
	return nil
}

func init() {
	cardImageCmd.AddCommand(cardImageDeleteCmd)
}
