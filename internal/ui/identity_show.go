package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayIdentity(w io.Writer, identity *fizzy.GetMyIdentityResponse) error {
	if len(identity.Accounts) == 0 {
		fmt.Fprintf(w, "No accounts found\n")
		return nil
	}

	// Display current user info from first account (user is same across accounts)
	user := identity.Accounts[0].User
	fmt.Fprintf(w, "User: %s\n", user.Name)
	fmt.Fprintf(w, "Email: %s\n", user.Email)
	fmt.Fprintf(w, "Role: %s\n", user.Role)
	fmt.Fprintf(w, "\nAvailable Accounts:\n")

	for i, account := range identity.Accounts {
		fmt.Fprintf(w, "  %d. %s %s\n", i+1, account.Name, DisplayMeta("slug", account.Slug))
	}

	return nil
}
