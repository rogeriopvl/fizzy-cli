package cmd

import "github.com/spf13/cobra"

var webhookDeliveryCmd = &cobra.Command{
	Use:   "delivery",
	Short: "Manage webhook deliveries",
	Long:  `View delivery attempts for a webhook`,
}

func init() {
	webhookCmd.AddCommand(webhookDeliveryCmd)
}
