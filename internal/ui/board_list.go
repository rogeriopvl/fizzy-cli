package ui

import (
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/api"
)

func DisplayBoards(boards []api.Board) error {
	for _, board := range boards {
		fmt.Printf("%s (%s)\n", board.Name, DisplayID(board.ID))
	}
	return nil
}
