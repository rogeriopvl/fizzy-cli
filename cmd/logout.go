package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and destroy the session",
	Long:  `Destroy the server-side session and log out the current user`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleLogout(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleLogout(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if err := a.Client.DeleteSession(context.Background()); err != nil {
		return fmt.Errorf("logging out: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully logged out\n")
	return nil
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
