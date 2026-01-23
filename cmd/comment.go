// Package cmd
package cmd

import "github.com/spf13/cobra"

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage card comments",
	Long:  `Manage comments on cards in Fizzy`,
}

func init() {
	rootCmd.AddCommand(commentCmd)
}
