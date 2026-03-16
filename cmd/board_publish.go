package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var boardPublishCmd = &cobra.Command{
	Use:   "publish <board_id>",
	Short: "Publish a board",
	Long:  `Make a board publicly accessible`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handlePublishBoard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handlePublishBoard(cmd *cobra.Command, boardID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	_, err := a.Client.PublishBoard(context.Background(), boardID)
	if err != nil {
		return fmt.Errorf("publishing board: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Board '%s' published successfully\n", boardID)
	return nil
}

func init() {
	boardCmd.AddCommand(boardPublishCmd)
}
