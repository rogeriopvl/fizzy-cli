package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
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

	// Read flag values directly from command
	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	status, _ := cmd.Flags().GetString("status")
	imageURL, _ := cmd.Flags().GetString("image-url")
	tagIDs, _ := cmd.Flags().GetStringSlice("tag-id")
	createdAt, _ := cmd.Flags().GetString("created-at")
	lastActiveAt, _ := cmd.Flags().GetString("last-active-at")

	payload := api.CreateCardPayload{
		Title:        title,
		Description:  description,
		Status:       status,
		ImageURL:     imageURL,
		TagIDS:       tagIDs,
		CreatedAt:    createdAt,
		LastActiveAt: lastActiveAt,
	}

	_, err := a.Client.PostCards(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("creating card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Card '%s' created successfully\n", title)
	return nil
}

func init() {
	cardCreateCmd.Flags().StringP("title", "t", "", "Card title (required)")
	cardCreateCmd.MarkFlagRequired("title")
	cardCreateCmd.Flags().StringP("description", "d", "", "Card description")
	cardCreateCmd.Flags().String("status", "", "Card status")
	cardCreateCmd.Flags().String("image-url", "", "Card image URL")
	cardCreateCmd.Flags().StringSlice("tag-id", []string{}, "Tag ID (can be used multiple times)")
	cardCreateCmd.Flags().String("created-at", "", "Creation timestamp")
	cardCreateCmd.Flags().String("last-active-at", "", "Last active timestamp")

	cardCmd.AddCommand(cardCreateCmd)
}
