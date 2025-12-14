package ui

import (
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/api"
)

func DisplayCards(cards []api.Card) error {
	for _, card := range cards {
		fmt.Printf("%d - %s\n", card.Number, card.Title)
	}
	return nil
}
