package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var accountListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  `Retrieve and display all accounts from Fizzy`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListAccounts(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListAccounts(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	identity, err := a.Client.GetMyIdentity(context.Background())
	if err != nil {
		return fmt.Errorf("fetching accounts: %w", err)
	}

	if len(identity.Accounts) == 0 {
		fmt.Println("No accounts found")
		return nil
	}

	return ui.DisplayAccounts(identity.Accounts)
}

func init() {
	accountCmd.AddCommand(accountListCmd)
}
