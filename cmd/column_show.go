package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var columnShowCmd = &cobra.Command{
	Use:   "show <column_id>",
	Short: "Show column details",
	Long:  `Retrieve and display detailed information about a specific column`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowColumnDetails(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowColumnDetails(cmd *cobra.Command, columnID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	column, err := a.Client.GetColumn(context.Background(), columnID)
	if err != nil {
		return fmt.Errorf("fetching column: %w", err)
	}

	return ui.DisplayColumn(cmd.OutOrStdout(), column)
}

func init() {
	columnCmd.AddCommand(columnShowCmd)
}
