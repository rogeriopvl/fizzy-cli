package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayExport(w io.Writer, e *fizzy.Export) error {
	fmt.Fprintf(w, "Export: %s\n", e.ID)
	fmt.Fprintf(w, "Status: %s\n", e.Status)
	fmt.Fprintf(w, "Created At: %s\n", FormatTime(e.CreatedAt))
	if e.DownloadURL != "" {
		fmt.Fprintf(w, "Download URL: %s\n", e.DownloadURL)
	}
	return nil
}
