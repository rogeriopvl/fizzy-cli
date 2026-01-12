package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardTagCmd = &cobra.Command{
	Use:   "tag <card_number> <tag_title>",
	Short: "Toggle a tag on or off for a card",
	Long: `Toggle a tag on or off for a card. If the tag doesn't exist, it will be created.

The tag title can be specified with or without a leading # symbol.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleTagCard(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleTagCard(cmd *cobra.Command, cardNumber, tagTitle string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	tagTitle = strings.TrimPrefix(tagTitle, "#")

	_, err = a.Client.PostCardTagging(context.Background(), cardNum, tagTitle)
	if err != nil {
		return fmt.Errorf("toggling tag on card: %w", err)
	}

	fmt.Printf("✓ Tag '%s' toggled on card #%d\n", tagTitle, cardNum)
	return nil
}

func init() {
	cardCmd.AddCommand(cardTagCmd)
}
