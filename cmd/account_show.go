package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var accountShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current account details",
	Long:  `Retrieve and display the current account's settings`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowAccount(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowAccount(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	account, err := a.Client.GetAccountSettings(context.Background())
	if err != nil {
		return fmt.Errorf("fetching account: %w", err)
	}

	return ui.DisplayAccount(cmd.OutOrStdout(), account)
}

func init() {
	accountCmd.AddCommand(accountShowCmd)
}
