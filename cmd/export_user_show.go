package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var exportUserShowCmd = &cobra.Command{
	Use:   "show <user_id> <export_id>",
	Short: "Show user data export status",
	Long:  `Retrieve and display the status of a user data export`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowUserDataExport(cmd, args[0], args[1]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowUserDataExport(cmd *cobra.Command, userID, exportID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	export, err := a.Client.GetUserDataExport(context.Background(), userID, exportID)
	if err != nil {
		return fmt.Errorf("fetching user data export: %w", err)
	}

	return ui.DisplayExport(cmd.OutOrStdout(), export)
}

func init() {
	exportUserCmd.AddCommand(exportUserShowCmd)
}
