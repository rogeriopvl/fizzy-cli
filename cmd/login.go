package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rogeriopvl/fizzy-cli/internal/api"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Prints instructions on how to authenticate with the Fizzy API",
	Long:  `Prints intructions on how to authenticate with the Fizzy API`,
	Run: func(cmd *cobra.Command, args []string) {
		token, isSet := os.LookupEnv("FIZZY_ACCESS_TOKEN")

		if isSet && token != "" {
			fmt.Printf("✓ Authenticated with access token: %s\n", token[:6]+"...")

			client, err := api.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating API client: %v\n", err)
				return
			}

			identity, err := client.GetMyIdentity(context.Background())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching identity: %v\n", err)
				return
			}

			fmt.Println("\nYour accounts:")
			for _, account := range identity.Accounts {
				fmt.Printf("  - %s\n", account.Name)
			}
		} else {
			fmt.Println("To authenticate with Fizzy's API you need an access token.")
			fmt.Println("\nGo to https://app.fizzy.do/<account_slug>/my/access_tokens and follow the instructions...")
			fmt.Println("(Replace <account_slug> with your account slug/id)")
			fmt.Println("\nThen export it as an environment variable in your shell, with the name FIZZY_ACCESS_TOKEN")
			fmt.Println("And re-run this command.")
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
