package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var notificationSettingsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show notification settings",
	Long:  `Retrieve and display the current user's notification settings`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowNotificationSettings(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowNotificationSettings(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	settings, err := a.Client.GetNotificationSettings(context.Background())
	if err != nil {
		return fmt.Errorf("fetching notification settings: %w", err)
	}

	return ui.DisplayNotificationSettings(cmd.OutOrStdout(), settings)
}

func init() {
	notificationSettingsCmd.AddCommand(notificationSettingsShowCmd)
}
