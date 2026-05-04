package cmd

import "github.com/spf13/cobra"

var exportAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage account exports",
	Long:  `Create and view account exports. Only admins and owners can create them.`,
}

func init() {
	exportCmd.AddCommand(exportAccountCmd)
}
