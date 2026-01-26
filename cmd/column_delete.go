package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var columnDeleteCmd = &cobra.Command{
	Use:   "delete <column_id>",
	Short: "Delete a column",
	Long:  `Delete a column. Only board administrators can delete columns.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteColumn(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteColumn(cmd *cobra.Command, columnID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err := a.Client.DeleteColumn(context.Background(), columnID)
	if err != nil {
		return fmt.Errorf("deleting column: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Column '%s' deleted successfully\n", columnID)
	return nil
}

func init() {
	columnCmd.AddCommand(columnDeleteCmd)
}
