package cmd

import "github.com/spf13/cobra"

var reactionCmd = &cobra.Command{
	Use:   "reaction",
	Short: "Manage comment reactions",
	Long:  `Manage reactions on comments in Fizzy`,
}

func init() {
	rootCmd.AddCommand(reactionCmd)
}
