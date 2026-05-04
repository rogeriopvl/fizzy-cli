package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayNotificationSettings(w io.Writer, settings *fizzy.NotificationSettings) error {
	fmt.Fprintf(w, "Bundle email frequency: %s\n", settings.BundleEmailFrequency)
	return nil
}
