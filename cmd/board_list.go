package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var boardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all boards",
	Long:  `Retrieve and display all boards from Fizzy`,
	RunE:  listBoards,
}

func listBoards(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// TODO: Implement API call to fetch boards
	boards, err := fetchBoards(ctx)
	if err != nil {
		return fmt.Errorf("failed to list boards: %w", err)
	}

	// TODO: Format and display boards
	if err := displayBoards(boards); err != nil {
		return fmt.Errorf("failed to display boards: %w", err)
	}

	return nil
}

// fetchBoards retrieves all boards from the Fizzy API.
func fetchBoards(ctx context.Context) ([]Board, error) {
	// TODO: Implement API call
	return nil, nil
}

// displayBoards formats and outputs the boards.
func displayBoards(boards []Board) error {
	// TODO: Implement output formatting
	for _, board := range boards {
		fmt.Printf("Board: %+v\n", board)
	}
	return nil
}

// Board represents a Fizzy board.
type Board struct {
	ID   string
	Name string
	// TODO: Add additional fields as needed
}

func init() {
	boardCmd.AddCommand(boardListCmd)
}
