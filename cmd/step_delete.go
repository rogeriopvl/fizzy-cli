package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var stepDeleteCmd = &cobra.Command{
	Use:   "delete <card_number> <step_id>",
	Short: "Delete a step",
	Long:  `Delete a step from a card`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteStep(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteStep(cmd *cobra.Command, cardNumber, stepID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.DeleteCardStep(context.Background(), cardNum, stepID)
	if err != nil {
		return fmt.Errorf("deleting step: %w", err)
	}

	fmt.Printf("✓ Step deleted successfully\n")
	return nil
}

func init() {
	stepCmd.AddCommand(stepDeleteCmd)
}
