package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var (
	updateTitle        string
	updateDescription  string
	updateStatus       string
	updateTagIDs       []string
	updateLastActiveAt string
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

	// Validate that at least one field is provided for update
	if updateTitle == "" && updateDescription == "" && updateStatus == "" && len(updateTagIDs) == 0 && updateLastActiveAt == "" {
		return fmt.Errorf("must provide at least one flag to update (--title, --description, --status, --tag-id, or --last-active-at)")
	}

	payload := api.UpdateCardPayload{
		Title:        updateTitle,
		Description:  updateDescription,
		Status:       updateStatus,
		TagIDS:       updateTagIDs,
		LastActiveAt: updateLastActiveAt,
	}

	card, err := a.Client.PutCard(context.Background(), cardNum, payload)
	if err != nil {
		return fmt.Errorf("updating card: %w", err)
	}

	fmt.Printf("✓ Card #%d updated successfully\n", card.Number)
	return nil
}

func init() {
	cardUpdateCmd.Flags().StringVarP(&updateTitle, "title", "t", "", "Card title")
	cardUpdateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "Card description")
	cardUpdateCmd.Flags().StringVar(&updateStatus, "status", "", "Card status")
	cardUpdateCmd.Flags().StringSliceVar(&updateTagIDs, "tag-id", []string{}, "Tag ID (can be used multiple times)")
	cardUpdateCmd.Flags().StringVar(&updateLastActiveAt, "last-active-at", "", "Last active timestamp")

	cardCmd.AddCommand(cardUpdateCmd)
}
