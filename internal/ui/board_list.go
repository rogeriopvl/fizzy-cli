package ui

import (
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayBoards(boards []fizzy.Board) error {
	for _, board := range boards {
		fmt.Printf("%s (%s)\n", board.Name, DisplayID(board.ID))
	}
	return nil
}
