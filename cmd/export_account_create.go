package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var exportAccountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Start an account export",
	Long:  `Start an account export job. Poll its status with 'export account show <id>'.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateAccountExport(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateAccountExport(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	export, err := a.Client.CreateAccountExport(context.Background())
	if err != nil {
		return fmt.Errorf("creating account export: %w", err)
	}

	return ui.DisplayExport(cmd.OutOrStdout(), export)
}

func init() {
	exportAccountCmd.AddCommand(exportAccountCreateCmd)
}
