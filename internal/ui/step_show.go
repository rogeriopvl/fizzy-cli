package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayStep(w io.Writer, step *fizzy.Step) error {
	check := "[ ]"
	if step.Completed {
		check = "[x]"
	}
	fmt.Fprintf(w, "%s %s (%s)\n", check, step.Content, DisplayID(step.ID))
	return nil
}
