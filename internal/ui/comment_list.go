package ui

import (
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayComments(comments []fizzy.Comment) error {
	for _, comment := range comments {
		fmt.Printf("%s - %s (%s)\n", comment.Creator.Name, comment.Body.PlainText, DisplayID(comment.ID))
	}
	return nil
}

func DisplayComment(comment *fizzy.Comment) error {
	fmt.Printf("Author: %s\n", comment.Creator.Name)
	fmt.Printf("Created: %s\n", comment.CreatedAt)
	if comment.UpdatedAt != comment.CreatedAt {
		fmt.Printf("Updated: %s\n", comment.UpdatedAt)
	}
	fmt.Printf("Card: %s\n", comment.Card.Title)
	fmt.Printf("\n%s\n", comment.Body.PlainText)
	return nil
}
