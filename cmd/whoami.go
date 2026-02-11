package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display current user identity and accounts",
	Long:  `Show information about the currently authenticated user and their available accounts`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleWhoami(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleWhoami(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	identity, err := a.Client.GetMyIdentity(context.Background())
	if err != nil {
		return fmt.Errorf("fetching identity: %w", err)
	}

	return ui.DisplayIdentity(cmd.OutOrStdout(), identity)
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
