package cmd

import "github.com/spf13/cobra"

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Manage exports",
	Long:  `Create and view account and user data exports`,
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
