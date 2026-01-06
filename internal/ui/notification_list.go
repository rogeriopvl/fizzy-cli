package ui

import (
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
)

func DisplayNotifications(notifications []api.Notification) error {
	for _, notification := range notifications {
		fmt.Printf("%s (%s)\n", notification.Title, DisplayID(notification.ID))
	}
	return nil
}
