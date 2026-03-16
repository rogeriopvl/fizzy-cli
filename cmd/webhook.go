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
  fizzy webhook list      List all webhooks`,
}

func init() {
	rootCmd.AddCommand(webhookCmd)
}
