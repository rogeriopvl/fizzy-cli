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

func DisplayComment(comment *api.Comment) error {
	fmt.Printf("Author: %s\n", comment.Creator.Name)
	fmt.Printf("Created: %s\n", comment.CreatedAt)
	if comment.UpdatedAt != comment.CreatedAt {
		fmt.Printf("Updated: %s\n", comment.UpdatedAt)
	}
	fmt.Printf("Card: %s\n", comment.Card.Title)
	fmt.Printf("\n%s\n", comment.Body.PlainText)
	return nil
}
