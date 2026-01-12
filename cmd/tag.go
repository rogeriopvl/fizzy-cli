package cmd

import (
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
	Long:  `Manage tags in Fizzy`,
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
