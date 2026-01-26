package ui

import (
	"fmt"
	"io"

	"github.com/rogeriopvl/fizzy/internal/api"
)

func DisplayBoard(w io.Writer, board *api.Board) error {
	fmt.Fprintf(w, "Board: %s\n", board.Name)
	fmt.Fprintf(w, "ID: %s\n", board.ID)
	fmt.Fprintf(w, "All Access: %v\n", board.AllAccess)
	fmt.Fprintf(w, "Created At: %s\n", FormatTime(board.CreatedAt))
	fmt.Fprintf(w, "Created By: %s\n", board.Creator.Name)
	return nil
}
