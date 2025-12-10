package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fizzy-cli",
	Short: "Fizzy CLI",
	Long:  `Fizzy CLI`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fizzy-cli.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
