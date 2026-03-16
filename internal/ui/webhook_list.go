package ui

import (
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayWebhooks(webhooks []fizzy.Webhook) error {
	for _, webhook := range webhooks {
		status := "inactive"
		if webhook.Active {
			status = "active"
		}
		fmt.Printf("%s (%s) [%s]\n", webhook.Name, DisplayID(webhook.ID), status)
	}
	return nil
}
