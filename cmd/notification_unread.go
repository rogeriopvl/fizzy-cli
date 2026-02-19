package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var notificationUnreadCmd = &cobra.Command{
	Use:   "unread <notification_id>",
	Short: "Mark notification as unread",
	Long:  `Mark a notification as unread`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUnreadNotification(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUnreadNotification(cmd *cobra.Command, notificationID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err := a.Client.MarkNotificationUnread(context.Background(), notificationID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("notification not found")
		}
		return fmt.Errorf("marking notification as unread: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Notification marked as unread successfully\n")
	return nil
}

func init() {
	notificationCmd.AddCommand(notificationUnreadCmd)
}
