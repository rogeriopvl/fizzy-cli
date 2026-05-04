package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var accountJoincodeResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the account join code",
	Long:  `Generate a new join code, invalidating the existing one. Requires admin role.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleResetJoinCode(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleResetJoinCode(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if err := a.Client.ResetAccountJoinCode(context.Background()); err != nil {
		return fmt.Errorf("resetting join code: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Join code reset; the previous code is no longer valid\n")
	return nil
}

func init() {
	accountJoincodeCmd.AddCommand(accountJoincodeResetCmd)
}
