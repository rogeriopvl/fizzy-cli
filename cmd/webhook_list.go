package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var webhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all webhooks",
	Long:  `Retrieve and display all webhooks for a board`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListWebhooks(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListWebhooks(cmd *cobra.Command) error {
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

	opts := &fizzy.ListOptions{}
	if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
		opts.Limit = limit
	}

	webhooks, err := a.Client.GetWebhooks(context.Background(), boardID, opts)
	if err != nil {
		return fmt.Errorf("fetching webhooks: %w", err)
	}

	if len(webhooks) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No webhooks found")
		return nil
	}

	return ui.DisplayWebhooks(webhooks)
}

func init() {
	webhookListCmd.Flags().StringP("board-id", "b", "", "Board ID (uses selected board if not specified)")
	webhookListCmd.Flags().IntP("limit", "l", 0, "Maximum number of webhooks to return (0 = no limit)")

	webhookCmd.AddCommand(webhookListCmd)
}
