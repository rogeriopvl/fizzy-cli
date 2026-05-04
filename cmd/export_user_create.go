package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var exportUserCreateCmd = &cobra.Command{
	Use:   "create <user_id>",
	Short: "Start a user data export",
	Long:  `Start a personal data export for the given user. You can only export data for your own user record.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateUserDataExport(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateUserDataExport(cmd *cobra.Command, userID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	export, err := a.Client.CreateUserDataExport(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("creating user data export: %w", err)
	}

	return ui.DisplayExport(cmd.OutOrStdout(), export)
}

func init() {
	exportUserCmd.AddCommand(exportUserCreateCmd)
}
