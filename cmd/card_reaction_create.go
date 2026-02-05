package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardReactionCreateCmd = &cobra.Command{
	Use:   "create <card_number> <emoji>",
	Short: "Create a reaction on a card",
	Long:  `Create an emoji reaction (boost) on a card`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateCardReaction(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateCardReaction(cmd *cobra.Command, cardNumber, emoji string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err = a.Client.PostCardReaction(context.Background(), cardNum, emoji)
	if err != nil {
		return fmt.Errorf("creating reaction: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Reaction %s created successfully\n", emoji)
	return nil
}

func init() {
	cardReactionCmd.AddCommand(cardReactionCreateCmd)
}
