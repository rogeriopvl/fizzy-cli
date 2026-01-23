package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var commentListCmd = &cobra.Command{
	Use:   "list <card_number>",
	Short: "List comments on a card",
	Long:  `Retrieve and display all comments on a card`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListComments(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListComments(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	comments, err := a.Client.GetCardComments(context.Background(), cardNum)
	if err != nil {
		return fmt.Errorf("fetching comments: %w", err)
	}

	if len(comments) == 0 {
		fmt.Println("No comments found")
		return nil
	}

	return ui.DisplayComments(comments)
}

func init() {
	commentCmd.AddCommand(commentListCmd)
}
