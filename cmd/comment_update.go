package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var commentUpdateCmd = &cobra.Command{
	Use:   "update <card_number> <comment_id>",
	Short: "Update an existing comment",
	Long:  `Update the body of an existing comment on a card`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateComment(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateComment(cmd *cobra.Command, cardNumber, commentID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	body, _ := cmd.Flags().GetString("body")

	comment, err := a.Client.UpdateCardComment(context.Background(), cardNum, commentID, body)
	if err != nil {
		return fmt.Errorf("updating comment: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Comment updated successfully (id: %s)\n", comment.ID)
	return nil
}

func init() {
	commentUpdateCmd.Flags().StringP("body", "b", "", "New comment body (required)")
	commentUpdateCmd.MarkFlagRequired("body")

	commentCmd.AddCommand(commentUpdateCmd)
}
