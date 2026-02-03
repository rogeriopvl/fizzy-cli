package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var userDeactivateCmd = &cobra.Command{
	Use:   "deactivate <user_id>",
	Short: "Deactivate a user",
	Long:  `Deactivate a user. Only account administrators can deactivate users.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeactivateUser(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeactivateUser(cmd *cobra.Command, userID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err := a.Client.DeleteUser(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("deactivating user: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ User '%s' deactivated successfully\n", userID)
	return nil
}

func init() {
	userCmd.AddCommand(userDeactivateCmd)
}
