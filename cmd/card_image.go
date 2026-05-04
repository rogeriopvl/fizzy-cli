package cmd

import "github.com/spf13/cobra"

var cardImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Manage card images",
	Long:  `Manage the image attached to a card`,
}

func init() {
	cardCmd.AddCommand(cardImageCmd)
}
