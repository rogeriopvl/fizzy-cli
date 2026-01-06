package cmd

import (
	"context"
	"fmt"

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

	notifications, err := a.Client.GetNotifications(context.Background())
	if err != nil {
		return fmt.Errorf("fetching notifications: %w", err)
	}

	if len(notifications) == 0 {
		fmt.Println("No notifications found")
		return nil
	}

	return ui.DisplayNotifications(notifications)
}

func init() {
	notificationCmd.AddCommand(notificationListCmd)
}
