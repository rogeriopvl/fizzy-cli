package ui

import (
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
)

func DisplayAccounts(accounts []api.Account) error {
	for _, account := range accounts {
		fmt.Printf("%s (%s)\n", account.Name, DisplayMeta("slug", account.Slug))
	}
	return nil
}
