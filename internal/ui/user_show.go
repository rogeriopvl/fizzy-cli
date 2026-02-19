package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayUser(w io.Writer, user *fizzy.User) error {
	status := "active"
	if !user.Active {
		status = "inactive"
	}

	fmt.Fprintf(w, "User: %s\n", user.Name)
	fmt.Fprintf(w, "ID: %s\n", user.ID)
	fmt.Fprintf(w, "Email: %s\n", user.Email)
	fmt.Fprintf(w, "Role: %s\n", user.Role)
	fmt.Fprintf(w, "Status: %s\n", status)
	fmt.Fprintf(w, "Created At: %s\n", FormatTime(user.CreatedAt))
	return nil
}
