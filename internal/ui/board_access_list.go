package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayBoardAccesses(w io.Writer, accesses *fizzy.BoardAccesses) error {
	fmt.Fprintf(w, "Board: %s (all-access: %v)\n", accesses.BoardID, accesses.AllAccess)
	for _, u := range accesses.Users {
		access := "no"
		if u.HasAccess {
			access = "yes"
		}
		involvement := u.Involvement
		if involvement == "" {
			involvement = "-"
		}
		fmt.Fprintf(w, "  %s (%s) access=%s involvement=%s\n", u.Name, DisplayID(u.ID), access, involvement)
	}
	return nil
}
