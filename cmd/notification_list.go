package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var notificationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notifications",
	Long:  `Retrieve and display all notifications from Fizzy`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListNotifications(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListNotifications(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	opts := &api.ListOptions{}
	if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
		opts.Limit = limit
	}

	notifications, err := a.Client.GetNotifications(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("fetching notifications: %w", err)
	}

	read, _ := cmd.Flags().GetBool("read")
	unread, _ := cmd.Flags().GetBool("unread")

	filtered := filterNotifications(notifications, read, unread)

	if len(filtered) == 0 {
		fmt.Println("No notifications found")
		return nil
	}

	return ui.DisplayNotifications(filtered)
}

func filterNotifications(notifications []api.Notification, read bool, unread bool) []api.Notification {
	if !read && !unread {
		return notifications
	}

	var filtered []api.Notification
	for _, notification := range notifications {
		if read && notification.Read {
			filtered = append(filtered, notification)
		} else if unread && !notification.Read {
			filtered = append(filtered, notification)
		}
	}

	return filtered
}

func init() {
	notificationCmd.AddCommand(notificationListCmd)
	notificationListCmd.Flags().BoolP("read", "r", false, "Show only read notifications")
	notificationListCmd.Flags().BoolP("unread", "u", false, "Show only unread notifications")
	notificationListCmd.Flags().IntP("limit", "l", 0, "Maximum number of notifications to return (0 = no limit)")
}
