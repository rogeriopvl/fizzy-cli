// Package cmd
package cmd

import "github.com/spf13/cobra"

var columnCmd = &cobra.Command{
	Use:   "column",
	Short: "Manage columns",
	Long:  `Manage columns in Fizzy`,
}

func init() {
	rootCmd.AddCommand(columnCmd)
}
