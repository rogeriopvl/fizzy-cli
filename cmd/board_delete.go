package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var boardDeleteCmd = &cobra.Command{
	Use:   "delete <board_id>",
	Short: "Delete a board",
	Long:  `Delete a board. Only board administrators can delete boards.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteBoard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteBoard(cmd *cobra.Command, boardID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err := a.Client.DeleteBoard(context.Background(), boardID)
	if err != nil {
		return fmt.Errorf("deleting board: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Board '%s' deleted successfully\n", boardID)
	return nil
}

func init() {
	boardCmd.AddCommand(boardDeleteCmd)
}
