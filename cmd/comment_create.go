package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var commentCreateCmd = &cobra.Command{
	Use:   "create <card_number>",
	Short: "Create a new comment",
	Long:  `Create a new comment on a card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateComment(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateComment(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	body, _ := cmd.Flags().GetString("body")

	_, err = a.Client.CreateCardComment(context.Background(), cardNum, body)
	if err != nil {
		return fmt.Errorf("creating comment: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Comment created successfully\n")
	return nil
}

func init() {
	commentCreateCmd.Flags().StringP("body", "b", "", "Comment body (required)")
	commentCreateCmd.MarkFlagRequired("body")

	commentCmd.AddCommand(commentCreateCmd)
}
