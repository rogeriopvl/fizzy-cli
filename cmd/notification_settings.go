package cmd

import "github.com/spf13/cobra"

var notificationSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Manage notification settings",
	Long:  `View and update notification settings for the current user`,
}

func init() {
	notificationCmd.AddCommand(notificationSettingsCmd)
}
