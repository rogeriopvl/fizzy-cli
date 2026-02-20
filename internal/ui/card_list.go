package ui

import (
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayCards(cards []fizzy.Card) error {
	for _, card := range cards {
		fmt.Printf("%d - %s\n", card.Number, card.Title)
	}
	return nil
}
