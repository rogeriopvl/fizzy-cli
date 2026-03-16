package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var webhookUpdateCmd = &cobra.Command{
	Use:   "update <webhook_id>",
	Short: "Update a webhook",
	Long:  `Update webhook settings such as name and subscribed actions`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateWebhook(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateWebhook(cmd *cobra.Command, webhookID string) error {
	if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("actions") {
		return fmt.Errorf("at least one flag must be provided (--name or --actions)")
	}

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

	payload := fizzy.UpdateWebhookPayload{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		payload.Name = name
	}
	if cmd.Flags().Changed("actions") {
		actions, _ := cmd.Flags().GetStringSlice("actions")
		payload.SubscribedActions = actions
	}

	_, err := a.Client.UpdateWebhook(context.Background(), boardID, webhookID, payload)
	if err != nil {
		return fmt.Errorf("updating webhook: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Webhook '%s' updated successfully\n", webhookID)
	return nil
}

func init() {
	webhookUpdateCmd.Flags().StringP("board-id", "b", "", "Board ID (uses selected board if not specified)")
	webhookUpdateCmd.Flags().StringP("name", "n", "", "Webhook name")
	webhookUpdateCmd.Flags().StringSliceP("actions", "a", nil, "Subscribed actions (comma-separated)")

	webhookCmd.AddCommand(webhookUpdateCmd)
}
