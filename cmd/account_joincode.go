package cmd

import "github.com/spf13/cobra"

var accountJoincodeCmd = &cobra.Command{
	Use:   "joincode",
	Short: "Manage the account join code",
	Long:  `Manage the account's join code (used to invite new users)`,
}

func init() {
	accountCmd.AddCommand(accountJoincodeCmd)
}
