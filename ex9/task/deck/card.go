package deck

import (
	"math/rand"
	"sort"
)

type Card struct {
	CardSuite string
	CardIndex int // 0 -> A, 1, 2, 3 ..., 10 -> J, 11 -> Q, 12 -> K
}

var SUITE_VAL map[string]int = map[string]int{
	"Spade":    0,
	"Diamonds": 1,
	"Club":     2,
	"Hearts":   3,
}
var SUITE_NAME []string = []string{"Spade", "Diamonds", "Club", "Hearts"}

type DeckOfCards struct {
	CardDeck     []Card
	lessCallback func(i, j int) bool
}

func (d *DeckOfCards) Len() int {
	return len(d.CardDeck)
}

func (d *DeckOfCards) Swap(i, j int) {
	if !(i >= 0 && j >= 0 && i <= len(d.CardDeck) && j <= len(d.CardDeck)) {
		return
	}
	d.CardDeck[i], d.CardDeck[j] = d.CardDeck[j], d.CardDeck[i]
}

func (d *DeckOfCards) Less(i, j int) bool {
	if d.lessCallback == nil {
		return d.LessDefault(i, j)
	}

	return d.lessCallback(i, j)
}

func (d *DeckOfCards) LessDefault(i, j int) bool {
	if !(i >= 0 && j >= 0 && i <= len(d.CardDeck) && j <= len(d.CardDeck)) {
		return false
	}
	suiteI, suiteJ := d.CardDeck[i].CardSuite, d.CardDeck[j].CardSuite
	if suiteI == suiteJ {
		return d.CardDeck[i].CardIndex < d.CardDeck[j].CardIndex
	}
	vi, ok1 := SUITE_VAL[suiteI]
	vj, ok2 := SUITE_VAL[suiteJ]
	if !ok1 || !ok2 {
		return false
	}

	return vi < vj
}

func (d *DeckOfCards) SortCustom(less func(i, j int) bool) *DeckOfCards {
	d.lessCallback = less
	sort.Sort(d)
	d.lessCallback = nil
	return d
}

func (d *DeckOfCards) Sort() *DeckOfCards {
	sort.Sort(d)
	return d
}

func (d *DeckOfCards) Shuffle() *DeckOfCards {
	rand.Shuffle(d.Len(), func(i, j int) {
		d.CardDeck[i], d.CardDeck[j] = d.CardDeck[j], d.CardDeck[i]
	})
	return d
}

func (d *DeckOfCards) AddJokers(rangeMin, rangeMax int) *DeckOfCards {
	if rangeMax-rangeMin <= 0 {
		return d
	}
	countToAdd := rangeMin + rand.Intn(rangeMax-rangeMin+1)

	//cardsToAdd := make([]Card, 0)

	for _ = range countToAdd {
		cardSuite := rand.Intn(4)
		c := Card{
			CardSuite: SUITE_NAME[cardSuite],
			CardIndex: 10,
		}

		d.CardDeck = append(d.CardDeck, c)
	}
	return d

}

func (d *DeckOfCards) Filter(hit func(c Card) bool) *DeckOfCards {
	toRemove := make([]int, 0)
	for idx, card := range d.CardDeck {
		if hit(card) {
			toRemove = append(toRemove, idx)
		}
	}

	for _, idx := range toRemove {
		d.CardDeck = append(d.CardDeck[:idx], d.CardDeck[idx+1:]...)
	}
	return d
}

func (d *DeckOfCards) MultiDeck(count int) *DeckOfCards {
	for _ = range count {
		for _, suite := range []string{"Spade", "Club", "Hearts", "Diamonds"} {
			for idx := range 13 {
				// for idx := 0; idx < 13; idx ++ {}
				d.CardDeck = append(d.CardDeck, Card{CardSuite: suite, CardIndex: idx})
			}
		}
	}
	return d
}

func New() *DeckOfCards {
	d := DeckOfCards{}
	for _, suite := range []string{"Spade", "Club", "Hearts", "Diamonds"} {
		for idx := range 13 {
			// for idx := 0; idx < 13; idx ++ {}
			d.CardDeck = append(d.CardDeck, Card{CardSuite: suite, CardIndex: idx})
		}
	}

	return &d
}
