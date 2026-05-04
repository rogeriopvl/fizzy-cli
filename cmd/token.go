package cmd

import "github.com/spf13/cobra"

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage personal access tokens",
	Long:  `Create personal access tokens for the Fizzy API`,
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
