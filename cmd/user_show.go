package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var userShowCmd = &cobra.Command{
	Use:   "show <user_id>",
	Short: "Show user details",
	Long:  `Retrieve and display detailed information about a specific user`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowUser(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowUser(cmd *cobra.Command, userID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	user, err := a.Client.GetUser(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("fetching user: %w", err)
	}

	return ui.DisplayUser(cmd.OutOrStdout(), user)
}

func init() {
	userCmd.AddCommand(userShowCmd)
}
