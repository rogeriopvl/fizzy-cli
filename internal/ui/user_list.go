package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayUsers(w io.Writer, users []fizzy.User) error {
	for _, user := range users {
		status := "active"
		if !user.Active {
			status = "inactive"
		}
		fmt.Fprintf(w, "%s (%s) - %s [%s]\n", user.Name, user.Email, DisplayID(user.ID), status)
	}
	return nil
}
