package cmd

import "github.com/spf13/cobra"

var boardAccessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage board access",
	Long:  `View user access for a board`,
}

func init() {
	boardCmd.AddCommand(boardAccessCmd)
}
