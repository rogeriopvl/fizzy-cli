package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var notificationReadAllCmd = &cobra.Command{
	Use:   "read-all",
	Short: "Mark all unread notifications as read",
	Long:  `Mark all unread notifications as read`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleReadAllNotifications(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleReadAllNotifications(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if err := a.Client.MarkAllNotificationsRead(context.Background()); err != nil {
		return fmt.Errorf("marking all notifications as read: %w", err)
	}

	fmt.Println("✓ All notifications marked as read successfully")
	return nil
}

func init() {
	notificationCmd.AddCommand(notificationReadAllCmd)
}
