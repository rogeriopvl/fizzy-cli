package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var notificationReadCmd = &cobra.Command{
	Use:   "read <notification_id>",
	Short: "Mark notification as read and display it",
	Long:  `Mark a notification as read and display its content`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleReadNotification(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleReadNotification(cmd *cobra.Command, notificationID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if err := a.Client.MarkNotificationRead(context.Background(), notificationID); err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("notification not found")
		}
		return fmt.Errorf("marking notification as read: %w", err)
	}

	notification, err := a.Client.GetNotification(context.Background(), notificationID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("notification not found")
		}
		return fmt.Errorf("fetching notification: %w", err)
	}

	return ui.DisplayNotification(notification)
}

func init() {
	notificationCmd.AddCommand(notificationReadCmd)
}
