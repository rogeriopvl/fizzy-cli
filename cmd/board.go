package cmd

import "github.com/spf13/cobra"

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Manage boards",
	Long:  `Manage boards in Fizzy`,
}

func init() {
	rootCmd.AddCommand(boardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// boardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// boardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
