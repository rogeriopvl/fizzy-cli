package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var commentDeleteCmd = &cobra.Command{
	Use:   "delete <card_number> <comment_id>",
	Short: "Delete a comment",
	Long:  `Delete a comment from a card`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteComment(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteComment(cmd *cobra.Command, cardNumber, commentID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.DeleteCardComment(context.Background(), cardNum, commentID)
	if err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Comment deleted successfully\n")
	return nil
}

func init() {
	commentCmd.AddCommand(commentDeleteCmd)
}
