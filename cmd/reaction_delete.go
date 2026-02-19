package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var reactionDeleteCmd = &cobra.Command{
	Use:   "delete <card_number> <comment_id> <reaction_id>",
	Short: "Delete a reaction from a comment",
	Long:  `Delete an emoji reaction from a comment`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteReaction(cmd, args[0], args[1], args[2]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteReaction(cmd *cobra.Command, cardNumber, commentID, reactionID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err = a.Client.DeleteCommentReaction(context.Background(), cardNum, commentID, reactionID)
	if err != nil {
		return fmt.Errorf("deleting reaction: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Reaction deleted successfully\n")
	return nil
}

func init() {
	reactionCmd.AddCommand(reactionDeleteCmd)
}
