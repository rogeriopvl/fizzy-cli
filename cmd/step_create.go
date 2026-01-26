package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var stepCreateCmd = &cobra.Command{
	Use:   "create <card_number>",
	Short: "Create a new step",
	Long:  `Create a new step (to-do item) on a card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateStep(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateStep(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	content, _ := cmd.Flags().GetString("content")
	completed, _ := cmd.Flags().GetBool("completed")

	_, err = a.Client.PostCardStep(context.Background(), cardNum, content, completed)
	if err != nil {
		return fmt.Errorf("creating step: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Step created successfully\n")
	return nil
}

func init() {
	stepCreateCmd.Flags().StringP("content", "c", "", "Step content (required)")
	stepCreateCmd.MarkFlagRequired("content")
	stepCreateCmd.Flags().BoolP("completed", "d", false, "Mark step as completed")

	stepCmd.AddCommand(stepCreateCmd)
}
