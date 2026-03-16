package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var webhookCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new webhook",
	Long:  `Create a new webhook for a board in Fizzy`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateWebhook(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateWebhook(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	boardID, _ := cmd.Flags().GetString("board-id")
	if boardID == "" {
		boardID = a.Config.SelectedBoard
	}
	if boardID == "" {
		return fmt.Errorf("no board specified: use --board-id or select a board with 'fizzy use'")
	}

	name, _ := cmd.Flags().GetString("name")
	url, _ := cmd.Flags().GetString("url")
	actions, _ := cmd.Flags().GetStringSlice("actions")

	payload := fizzy.CreateWebhookPayload{
		Name:              name,
		URL:               url,
		SubscribedActions: actions,
	}

	webhook, err := a.Client.CreateWebhook(context.Background(), boardID, payload)
	if err != nil {
		return fmt.Errorf("creating webhook: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Webhook '%s' created successfully (ID: %s)\n", webhook.Name, webhook.ID)
	return nil
}

func init() {
	webhookCreateCmd.Flags().StringP("board-id", "b", "", "Board ID (uses selected board if not specified)")
	webhookCreateCmd.Flags().StringP("name", "n", "", "Webhook name (required)")
	webhookCreateCmd.MarkFlagRequired("name")
	webhookCreateCmd.Flags().StringP("url", "u", "", "Webhook payload URL (required)")
	webhookCreateCmd.MarkFlagRequired("url")
	webhookCreateCmd.Flags().StringSliceP("actions", "a", nil, "Subscribed actions (comma-separated)")

	webhookCmd.AddCommand(webhookCreateCmd)
}
