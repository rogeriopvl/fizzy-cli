package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var stepShowCmd = &cobra.Command{
	Use:   "show <card_number> <step_id>",
	Short: "Show step details",
	Long:  `Retrieve and display a single step on a card`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowStep(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowStep(cmd *cobra.Command, cardNumber, stepID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	step, err := a.Client.GetCardStep(context.Background(), cardNum, stepID)
	if err != nil {
		return fmt.Errorf("fetching step: %w", err)
	}

	return ui.DisplayStep(cmd.OutOrStdout(), step)
}

func init() {
	stepCmd.AddCommand(stepShowCmd)
}
