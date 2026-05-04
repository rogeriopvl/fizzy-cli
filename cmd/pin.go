package cmd

import "github.com/spf13/cobra"

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Manage pinned cards",
	Long:  `Manage the current user's pinned cards`,
}

func init() {
	rootCmd.AddCommand(pinCmd)
}
