package cmd

import (
	"fmt"
	"os"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "fizzy-cli",
	Short:   "Fizzy CLI",
	Long:    `Fizzy CLI`,
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		a, _ := app.New(Version)
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
	rootCmd.SetVersionTemplate(fmt.Sprintf("fizzy-cli v%s\n", Version))
}
