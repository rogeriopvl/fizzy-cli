// Package cmd
package cmd

import "github.com/spf13/cobra"

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Manage notifications",
	Long:  `Manage notifications in Fizzy`,
}

func init() {
	rootCmd.AddCommand(notificationCmd)
}
