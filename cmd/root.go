package cmd

import (
	"os"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fizzy-cli",
	Short: "Fizzy CLI",
	Long:  `Fizzy CLI`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		a, _ := app.New()
		if a != nil {
			cmd.SetContext(a.ToContext(cmd.Context()))
		}
	},
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
