package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var reactionListCmd = &cobra.Command{
	Use:   "list <card_number> <comment_id>",
	Short: "List reactions on a comment",
	Long:  `Retrieve and display all reactions on a comment`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListReactions(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListReactions(cmd *cobra.Command, cardNumber, commentID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	reactions, err := a.Client.GetCommentReactions(context.Background(), cardNum, commentID)
	if err != nil {
		return fmt.Errorf("fetching reactions: %w", err)
	}

	if len(reactions) == 0 {
		fmt.Println("No reactions found")
		return nil
	}

	return ui.DisplayReactions(reactions)
}

func init() {
	reactionCmd.AddCommand(reactionListCmd)
}
