package cmd

import (
	"github.com/spf13/cobra"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage webhooks",
	Long: `Manage webhooks for a board.

Use subcommands to manage webhooks:
  fizzy webhook create    Create a new webhook
  fizzy webhook list      List all webhooks
  fizzy webhook show      Show webhook details
  fizzy webhook update    Update a webhook
  fizzy webhook delete    Delete a webhook
  fizzy webhook activate  Activate a webhook`,
}

func init() {
	rootCmd.AddCommand(webhookCmd)
}
