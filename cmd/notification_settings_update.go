package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var notificationSettingsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update notification settings",
	Long:  `Update the current user's notification settings`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateNotificationSettings(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateNotificationSettings(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("bundle-email-frequency") {
		return fmt.Errorf("at least one flag must be provided (--bundle-email-frequency)")
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	freq, _ := cmd.Flags().GetString("bundle-email-frequency")
	payload := fizzy.UpdateNotificationSettingsPayload{BundleEmailFrequency: freq}

	if err := a.Client.UpdateNotificationSettings(context.Background(), payload); err != nil {
		return fmt.Errorf("updating notification settings: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Notification settings updated\n")
	return nil
}

func init() {
	notificationSettingsUpdateCmd.Flags().String("bundle-email-frequency", "", "Bundle email frequency")
	notificationSettingsCmd.AddCommand(notificationSettingsUpdateCmd)
}
