package ui

import (
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
)

func DisplayReactions(reactions []api.Reaction) error {
	for _, reaction := range reactions {
		fmt.Printf("%s %s (%s)\n", reaction.Content, reaction.Reacter.Name, DisplayID(reaction.ID))
	}
	return nil
}
