package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var boardShowCmd = &cobra.Command{
	Use:   "show <board_id>",
	Short: "Show board details",
	Long:  `Retrieve and display detailed information about a specific board`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowBoardDetails(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowBoardDetails(cmd *cobra.Command, boardID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	board, err := a.Client.GetBoard(context.Background(), boardID)
	if err != nil {
		return fmt.Errorf("fetching board: %w", err)
	}

	return ui.DisplayBoard(cmd.OutOrStdout(), board)
}

func init() {
	boardCmd.AddCommand(boardShowCmd)
}
