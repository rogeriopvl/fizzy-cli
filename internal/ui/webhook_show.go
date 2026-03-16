package ui

import (
	"fmt"
	"io"
	"strings"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayWebhook(w io.Writer, webhook *fizzy.Webhook) error {
	fmt.Fprintf(w, "Webhook: %s\n", webhook.Name)
	fmt.Fprintf(w, "ID: %s\n", webhook.ID)
	fmt.Fprintf(w, "URL: %s\n", webhook.PayloadURL)
	fmt.Fprintf(w, "Active: %v\n", webhook.Active)
	if len(webhook.SubscribedActions) > 0 {
		fmt.Fprintf(w, "Actions: %s\n", strings.Join(webhook.SubscribedActions, ", "))
	}
	fmt.Fprintf(w, "Created At: %s\n", FormatTime(webhook.CreatedAt))
	return nil
}
