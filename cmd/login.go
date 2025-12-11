package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rogeriopvl/fizzy-cli/internal/api"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Prints instructions on how to authenticate with the Fizzy API",
	Long:  `Prints intructions on how to authenticate with the Fizzy API`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleLogin(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	},
}

func handleLogin(cmd *cobra.Command) error {
	token, isSet := os.LookupEnv("FIZZY_ACCESS_TOKEN")
	if !isSet || token == "" {
		return printAuthInstructions()
	}

	fmt.Printf("✓ Authenticated with access token: %s\n", token[:6]+"...")

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	identity, err := a.Client.GetMyIdentity(context.Background())
	if err != nil {
		return fmt.Errorf("fetching identity: %w", err)
	}

	selected, err := chooseAccount(identity.Accounts)
	if err != nil {
		return err
	}

	fmt.Printf("\nSelected account: %s (%s)\n", selected.Name, selected.Slug)

	// Save the selected account to config
	a.Config.SelectedAccount = selected.Slug
	if err := a.Config.Save(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	return nil
}

func chooseAccount(accounts []api.Account) (api.Account, error) {
	if len(accounts) == 1 {
		selected := accounts[0]
		fmt.Printf("\nUsing account: %s (%s)\n", selected.Name, selected.Slug)
		return selected, nil
	}

	fmt.Println("\nYour accounts:")
	return ui.SelectAccount(accounts)
}

func printAuthInstructions() error {
	fmt.Println("To authenticate with Fizzy's API you need an access token.")
	fmt.Printf("\nGo to https://app.fizzy.do/<account_slug>/my/access_tokens and follow the instructions...\n")
	fmt.Println("(Replace <account_slug> with your account slug)")
	fmt.Printf("\nThen export it as an environment variable in your shell, with the name FIZZY_ACCESS_TOKEN\n")
	fmt.Println("And re-run this command.")
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
