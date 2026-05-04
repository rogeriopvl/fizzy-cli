package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayActivities(w io.Writer, activities []fizzy.Activity) error {
	for _, a := range activities {
		fmt.Fprintf(w, "%s [%s] %s\n", FormatTime(a.CreatedAt), a.Action, a.Description)
	}
	return nil
}
