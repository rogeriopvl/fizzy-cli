package ui

import (
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
)

func DisplayComments(comments []api.Comment) error {
	for _, comment := range comments {
		fmt.Printf("%s - %s (%s)\n", comment.Creator.Name, comment.Body.PlainText, DisplayID(comment.ID))
	}
	return nil
}
