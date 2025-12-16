package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var (
	cardTitle        string
	cardDescription  string
	cardStatus       string
	cardImageURL     string
	cardTagIDs       []string
	cardCreatedAt    string
	cardLastActiveAt string
)

var cardCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new card",
	Long:  `Create a new card in the selected board`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateCard(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateCard(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if a.Config.SelectedBoard == "" {
		return fmt.Errorf("no board selected")
	}

	payload := api.CreateCardPayload{
		Title:        cardTitle,
		Description:  cardDescription,
		Status:       cardStatus,
		ImageURL:     cardImageURL,
		TagIDS:       cardTagIDs,
		CreatedAt:    cardCreatedAt,
		LastActiveAt: cardLastActiveAt,
	}

	_, err := a.Client.PostCards(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("creating card: %w", err)
	}

	fmt.Printf("✓ Card '%s' created successfully\n", cardTitle)
	return nil
}

func init() {
	cardCreateCmd.Flags().StringVarP(&cardTitle, "title", "t", "", "Card title (required)")
	cardCreateCmd.MarkFlagRequired("title")
	cardCreateCmd.Flags().StringVarP(&cardDescription, "description", "d", "", "Card description")
	cardCreateCmd.Flags().StringVar(&cardStatus, "status", "", "Card status")
	cardCreateCmd.Flags().StringVar(&cardImageURL, "image-url", "", "Card image URL")
	cardCreateCmd.Flags().StringSliceVar(&cardTagIDs, "tag-id", []string{}, "Tag ID (can be used multiple times)")
	cardCreateCmd.Flags().StringVar(&cardCreatedAt, "created-at", "", "Creation timestamp")
	cardCreateCmd.Flags().StringVar(&cardLastActiveAt, "last-active-at", "", "Last active timestamp")

	cardCmd.AddCommand(cardCreateCmd)
}
