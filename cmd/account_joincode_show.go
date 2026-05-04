package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var accountJoincodeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the account join code",
	Long:  `Retrieve and display the account's join code`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowJoinCode(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowJoinCode(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	jc, err := a.Client.GetAccountJoinCode(context.Background())
	if err != nil {
		return fmt.Errorf("fetching join code: %w", err)
	}

	return ui.DisplayJoinCode(cmd.OutOrStdout(), jc)
}

func init() {
	accountJoincodeCmd.AddCommand(accountJoincodeShowCmd)
}
