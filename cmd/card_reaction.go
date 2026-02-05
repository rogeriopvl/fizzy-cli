package cmd

import "github.com/spf13/cobra"

var cardReactionCmd = &cobra.Command{
	Use:   "reaction",
	Short: "Manage card reactions",
	Long:  `Manage reactions (boosts) on cards in Fizzy`,
}

func init() {
	cardCmd.AddCommand(cardReactionCmd)
}
