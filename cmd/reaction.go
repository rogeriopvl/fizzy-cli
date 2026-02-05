package cmd

import "github.com/spf13/cobra"

var reactionCmd = &cobra.Command{
	Use:   "reaction",
	Short: "Manage reactions",
	Long:  `Manage reactions on cards and comments in Fizzy`,
}

func init() {
	rootCmd.AddCommand(reactionCmd)
}
