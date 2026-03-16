package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var webhookDeleteCmd = &cobra.Command{
	Use:   "delete <webhook_id>",
	Short: "Delete a webhook",
	Long:  `Delete a webhook from a board`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteWebhook(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteWebhook(cmd *cobra.Command, webhookID string) error {
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

	err := a.Client.DeleteWebhook(context.Background(), boardID, webhookID)
	if err != nil {
		return fmt.Errorf("deleting webhook: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Webhook '%s' deleted successfully\n", webhookID)
	return nil
}

func init() {
	webhookDeleteCmd.Flags().StringP("board-id", "b", "", "Board ID (uses selected board if not specified)")

	webhookCmd.AddCommand(webhookDeleteCmd)
}
