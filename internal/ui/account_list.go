package ui

import (
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayAccounts(accounts []fizzy.Account) error {
	for _, account := range accounts {
		fmt.Printf("%s (%s)\n", account.Name, DisplayMeta("slug", account.Slug))
	}
	return nil
}
