package ui

import (
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayReactions(reactions []fizzy.Reaction) error {
	for _, reaction := range reactions {
		fmt.Printf("%s %s (%s)\n", reaction.Content, reaction.Reacter.Name, DisplayID(reaction.ID))
	}
	return nil
}
