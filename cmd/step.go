// Package cmd
package cmd

import "github.com/spf13/cobra"

var stepCmd = &cobra.Command{
	Use:   "step",
	Short: "Manage card steps",
	Long:  `Manage steps (to-do items) on cards in Fizzy`,
}

func init() {
	rootCmd.AddCommand(stepCmd)
}
