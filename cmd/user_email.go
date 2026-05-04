package cmd

import "github.com/spf13/cobra"

var userEmailCmd = &cobra.Command{
	Use:   "email",
	Short: "Manage user email addresses",
	Long:  `Request and confirm email address changes for users`,
}

func init() {
	userCmd.AddCommand(userEmailCmd)
}
