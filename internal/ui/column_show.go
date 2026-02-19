package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayColumn(w io.Writer, column *fizzy.Column) error {
	fmt.Fprintf(w, "Column: %s\n", column.Name)
	fmt.Fprintf(w, "ID: %s\n", column.ID)
	fmt.Fprintf(w, "Color: %s (%s)\n", column.Color.Name, column.Color.Value)
	fmt.Fprintf(w, "Created At: %s\n", FormatTime(column.CreatedAt))
	return nil
}
