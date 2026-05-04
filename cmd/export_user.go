package cmd

import "github.com/spf13/cobra"

var exportUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user data exports",
	Long:  `Create and view personal data exports. You can only manage exports for your own user record.`,
}

func init() {
	exportCmd.AddCommand(exportUserCmd)
}
