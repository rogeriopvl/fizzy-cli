package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var commentShowCmd = &cobra.Command{
	Use:   "show <card_number> <comment_id>",
	Short: "Show a specific comment",
	Long:  `Display details of a specific comment on a card`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowComment(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowComment(cmd *cobra.Command, cardNumber, commentID string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	comment, err := a.Client.GetCardComment(context.Background(), cardNum, commentID)
	if err != nil {
		return fmt.Errorf("fetching comment: %w", err)
	}

	return ui.DisplayComment(comment)
}

func init() {
	commentCmd.AddCommand(commentShowCmd)
}
