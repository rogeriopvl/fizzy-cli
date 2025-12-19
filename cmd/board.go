// Package cmd
package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Show the currently selected board",
	Long:  `Display the name and ID of the currently selected board.

Use subcommands to list, create, or manage boards:
  fizzy board list      List all boards
  fizzy board create    Create a new board`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleShowBoard(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleShowBoard(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if a.Config.SelectedBoard == "" {
		return fmt.Errorf("no board selected")
	}

	board, err := a.Client.GetBoard(context.Background(), a.Config.SelectedBoard)
	if err != nil {
		return fmt.Errorf("fetching board: %w", err)
	}

	fmt.Printf("%s (%s)\n", board.Name, ui.DisplayID(board.ID))
	return nil
}

func init() {
	rootCmd.AddCommand(boardCmd)
}
