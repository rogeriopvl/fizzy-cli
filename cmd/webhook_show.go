package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var webhookShowCmd = &cobra.Command{
	Use:   "show <webhook_id>",
	Short: "Show webhook details",
	Long:  `Retrieve and display detailed information about a specific webhook`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowWebhook(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowWebhook(cmd *cobra.Command, webhookID string) error {
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

	webhook, err := a.Client.GetWebhook(context.Background(), boardID, webhookID)
	if err != nil {
		return fmt.Errorf("fetching webhook: %w", err)
	}

	return ui.DisplayWebhook(cmd.OutOrStdout(), webhook)
}

func init() {
	webhookShowCmd.Flags().StringP("board-id", "b", "", "Board ID (uses selected board if not specified)")

	webhookCmd.AddCommand(webhookShowCmd)
}
