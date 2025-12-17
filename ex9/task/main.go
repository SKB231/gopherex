package main

import (
	"fmt"

	"github.com/SKB231/gopherex/ex9/task/deck"
)

func main() {
	fmt.Println("Hello World! Grabbing a deck of cards!")
	cards := deck.New()
	fmt.Println(cards)
	fmt.Println("Shuffling")
	cards.Shuffle()
	fmt.Println(cards)
	fmt.Println("Sorting!")
	cards.Sort()
	fmt.Println(cards)

	cards.Filter(func(c deck.Card) bool {
		return c.CardIndex == 2
	})


	fmt.Println(cards)
}
