package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Prints instructions on how to authenticate with the Fizzy API",
	Long:  `Prints intructions on how to authenticate with the Fizzy API`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To authenticate with Fizzy's API you need an access token.")
		fmt.Println("\nGo to https://app.fizzy.do/<account_slug>/my/access_tokens and follow the instructions...")
		fmt.Println("(Replace <account_slug> with your account slug/id)")
		fmt.Println("\nThen export it as an environment variable in your shell, with the name FIZZY_ACCESS_TOKEN")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
