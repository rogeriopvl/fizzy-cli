// Package cmd
package cmd

import "github.com/spf13/cobra"

var cardCmd = &cobra.Command{
	Use:   "card",
	Short: "Manage cards",
	Long:  `Manage cards in Fizzy`,
}

func init() {
	rootCmd.AddCommand(cardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cardsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
