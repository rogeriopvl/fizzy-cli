package ui

import (
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
)

func DisplayNotifications(notifications []api.Notification) error {
	for _, notification := range notifications {
		status := "read"
		if !notification.Read {
			status = "unread"
		}
		fmt.Printf("[%s] %s - %s (from %s)\n", status, notification.Title, notification.Body, notification.Creator.Name)
	}
	return nil
}
