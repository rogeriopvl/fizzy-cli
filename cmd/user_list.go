package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `Retrieve and display all users from the current account`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListUsers(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListUsers(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	users, err := a.Client.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("fetching users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}

	return ui.DisplayUsers(cmd.OutOrStdout(), users)
}

func init() {
	userCmd.AddCommand(userListCmd)
}
