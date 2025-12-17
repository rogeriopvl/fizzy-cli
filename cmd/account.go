// Package cmd
package cmd

import "github.com/spf13/cobra"

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
	Long:  `Manage accounts in Fizzy`,
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
