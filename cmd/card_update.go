package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var cardUpdateCmd = &cobra.Command{
	Use:   "update <card_number>",
	Short: "Update a card",
	Long:  `Update an existing card's details`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateCard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateCard(cmd *cobra.Command, cardNumber string) error {
	cardNum, err := strconv.Atoi(cardNumber)
	if err != nil {
		return fmt.Errorf("invalid card number: %w", err)
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	// Build payload only with flags that were explicitly set
	var payload api.UpdateCardPayload
	hasChanges := false

	if cmd.Flags().Changed("title") {
		payload.Title, _ = cmd.Flags().GetString("title")
		hasChanges = true
	}
	if cmd.Flags().Changed("description") {
		payload.Description, _ = cmd.Flags().GetString("description")
		hasChanges = true
	}
	if cmd.Flags().Changed("status") {
		payload.Status, _ = cmd.Flags().GetString("status")
		hasChanges = true
	}
	if cmd.Flags().Changed("tag-id") {
		payload.TagIDS, _ = cmd.Flags().GetStringSlice("tag-id")
		hasChanges = true
	}
	if cmd.Flags().Changed("last-active-at") {
		payload.LastActiveAt, _ = cmd.Flags().GetString("last-active-at")
		hasChanges = true
	}

	if !hasChanges {
		return fmt.Errorf("must provide at least one flag to update (--title, --description, --status, --tag-id, or --last-active-at)")
	}

	card, err := a.Client.PutCard(context.Background(), cardNum, payload)
	if err != nil {
		return fmt.Errorf("updating card: %w", err)
	}

	fmt.Printf("✓ Card #%d updated successfully\n", card.Number)
	return nil
}

func init() {
	cardUpdateCmd.Flags().StringP("title", "t", "", "Card title")
	cardUpdateCmd.Flags().StringP("description", "d", "", "Card description")
	cardUpdateCmd.Flags().String("status", "", "Card status")
	cardUpdateCmd.Flags().StringSlice("tag-id", []string{}, "Tag ID (can be used multiple times)")
	cardUpdateCmd.Flags().String("last-active-at", "", "Last active timestamp")

	cardCmd.AddCommand(cardUpdateCmd)
}
