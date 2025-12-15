package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var boardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all boards",
	Long:  `Retrieve and display all boards from Fizzy`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListBoards(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListBoards(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	boards, err := a.Client.GetBoards(context.Background())
	if err != nil {
		return fmt.Errorf("fetching boards: %w", err)
	}

	if len(boards) == 0 {
		fmt.Println("No boards found")
		return nil
	}

	return ui.DisplayBoards(boards)
}

func init() {
	boardCmd.AddCommand(boardListCmd)
}
