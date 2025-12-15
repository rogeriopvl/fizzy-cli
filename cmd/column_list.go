package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var columnListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all columns",
	Long:  `Retrieve and display all columns in the selected board`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListColumns(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListColumns(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	columns, err := a.Client.GetColumns(context.Background())
	if err != nil {
		return fmt.Errorf("fetching columns: %w", err)
	}

	if len(columns) == 0 {
		fmt.Println("No columns found")
		return nil
	}

	return ui.DisplayColumns(columns)
}

func init() {
	columnCmd.AddCommand(columnListCmd)
}
