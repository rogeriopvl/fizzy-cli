// Package cmd
package cmd

import (
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long: `Manage users in your account.

Use subcommands to list, view, or manage users:
  fizzy user list         List all users
  fizzy user show <id>    Show user details
  fizzy user update <id>  Update user settings
  fizzy user deactivate <id>  Deactivate a user`,
}

func init() {
	rootCmd.AddCommand(userCmd)
}
