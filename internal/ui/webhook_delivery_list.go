package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayWebhookDeliveries(w io.Writer, deliveries []fizzy.WebhookDelivery) error {
	for _, d := range deliveries {
		code := 0
		if d.Response != nil {
			code = d.Response.Code
		}
		fmt.Fprintf(w, "%s [%s] %s response=%d (%s)\n",
			FormatTime(d.CreatedAt), d.State, d.Event.Action, code, DisplayID(d.ID))
	}
	return nil
}
