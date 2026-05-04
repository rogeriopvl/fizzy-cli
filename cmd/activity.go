package cmd

import "github.com/spf13/cobra"

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "View account activity",
	Long:  `View the account-wide activity feed`,
}

func init() {
	rootCmd.AddCommand(activityCmd)
}
