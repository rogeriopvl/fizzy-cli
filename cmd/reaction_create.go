package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var reactionCreateCmd = &cobra.Command{
	Use:   "create <card_number> <comment_id> <emoji>",
	Short: "Create a reaction on a comment",
	Long:  `Create an emoji reaction on a comment`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateReaction(cmd, args[0], args[1], args[2]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateReaction(cmd *cobra.Command, cardNumber, commentID, emoji string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.PostCommentReaction(context.Background(), cardNum, commentID, emoji)
	if err != nil {
		return fmt.Errorf("creating reaction: %w", err)
	}

	fmt.Printf("✓ Reaction %s created successfully\n", emoji)
	return nil
}

func init() {
	reactionCmd.AddCommand(reactionCreateCmd)
}
