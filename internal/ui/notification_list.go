package ui

import (
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayNotifications(notifications []fizzy.Notification) error {
	for _, notification := range notifications {
		fmt.Printf("%s (%s)\n", notification.Title, DisplayID(notification.ID))
	}
	return nil
}

func DisplayNotification(notification *fizzy.Notification) error {
	status := "read"
	if !notification.Read {
		status = "unread"
	}
	fmt.Printf("%s (%s)\n", notification.Title, DisplayID(notification.ID))
	fmt.Printf("Status: %s\n", status)
	fmt.Printf("Card: %s\n", notification.Card.Title)
	fmt.Printf("From: %s\n", notification.Creator.Name)
	fmt.Printf("Message: %s\n", notification.Body)
	return nil
}
