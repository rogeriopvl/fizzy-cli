package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var boardUnpublishCmd = &cobra.Command{
	Use:   "unpublish <board_id>",
	Short: "Unpublish a board",
	Long:  `Remove public access from a board`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUnpublishBoard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUnpublishBoard(cmd *cobra.Command, boardID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	err := a.Client.UnpublishBoard(context.Background(), boardID)
	if err != nil {
		return fmt.Errorf("unpublishing board: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Board '%s' unpublished successfully\n", boardID)
	return nil
}

func init() {
	boardCmd.AddCommand(boardUnpublishCmd)
}
