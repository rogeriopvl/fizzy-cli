package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayAccount(w io.Writer, account *fizzy.Account) error {
	fmt.Fprintf(w, "Account: %s\n", account.Name)
	fmt.Fprintf(w, "ID: %s\n", account.ID)
	if account.Slug != "" {
		fmt.Fprintf(w, "Slug: %s\n", account.Slug)
	}
	fmt.Fprintf(w, "Cards: %d\n", account.CardsCount)
	fmt.Fprintf(w, "Auto-postpone period: %d days\n", account.AutoPostponePeriodInDays)
	fmt.Fprintf(w, "Created At: %s\n", FormatTime(account.CreatedAt))
	return nil
}
