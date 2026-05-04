package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var webhookDeliveryListCmd = &cobra.Command{
	Use:   "list <webhook_id>",
	Short: "List webhook deliveries",
	Long:  `Retrieve and display delivery attempts for a webhook`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListWebhookDeliveries(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListWebhookDeliveries(cmd *cobra.Command, webhookID string) error {
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

	deliveries, err := a.Client.GetWebhookDeliveries(context.Background(), boardID, webhookID, opts)
	if err != nil {
		return fmt.Errorf("fetching webhook deliveries: %w", err)
	}

	if len(deliveries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No deliveries found")
		return nil
	}

	return ui.DisplayWebhookDeliveries(cmd.OutOrStdout(), deliveries)
}

func init() {
	webhookDeliveryListCmd.Flags().StringP("board-id", "b", "", "Board ID (uses selected board if not specified)")
	webhookDeliveryListCmd.Flags().IntP("limit", "l", 0, "Maximum number of deliveries to return (0 = no limit)")

	webhookDeliveryCmd.AddCommand(webhookDeliveryListCmd)
}
