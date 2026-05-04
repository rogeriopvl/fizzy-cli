package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var exportAccountShowCmd = &cobra.Command{
	Use:   "show <export_id>",
	Short: "Show account export status",
	Long:  `Retrieve and display the status of an account export`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowAccountExport(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowAccountExport(cmd *cobra.Command, exportID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	export, err := a.Client.GetAccountExport(context.Background(), exportID)
	if err != nil {
		return fmt.Errorf("fetching account export: %w", err)
	}

	return ui.DisplayExport(cmd.OutOrStdout(), export)
}

func init() {
	exportAccountCmd.AddCommand(exportAccountShowCmd)
}
