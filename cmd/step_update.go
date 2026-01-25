package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var stepUpdateCmd = &cobra.Command{
	Use:   "update <card_number> <step_id>",
	Short: "Update an existing step",
	Long:  `Update the content or completion status of an existing step on a card`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateStep(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateStep(cmd *cobra.Command, cardNumber, stepID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	var contentPtr *string
	var completedPtr *bool

	if cmd.Flags().Changed("content") {
		content, _ := cmd.Flags().GetString("content")
		contentPtr = &content
	}

	if cmd.Flags().Changed("completed") {
		completed, _ := cmd.Flags().GetBool("completed")
		completedPtr = &completed
	}

	if contentPtr == nil && completedPtr == nil {
		return fmt.Errorf("at least one of --content or --completed must be provided")
	}

	step, err := a.Client.PutCardStep(context.Background(), cardNum, stepID, contentPtr, completedPtr)
	if err != nil {
		return fmt.Errorf("updating step: %w", err)
	}

	fmt.Printf("✓ Step updated successfully (id: %s)\n", step.ID)
	return nil
}

func init() {
	stepUpdateCmd.Flags().StringP("content", "c", "", "New step content")
	stepUpdateCmd.Flags().BoolP("completed", "d", false, "Mark step as completed")

	stepCmd.AddCommand(stepUpdateCmd)
}
