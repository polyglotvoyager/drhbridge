package main

import (
    "slices"
    "math/rand/v2"
)

var ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
var suits = []string{"C", "D", "H", "S"}

type Card struct {
    Rank string
    Suit string
}

// return index in ranks slice
func (c Card) RankIndex() int {
    return slices.Index(ranks, c.Rank)
}

func CardBeats(a, b Card, ledSuit string, trump string) bool {
    // cards are the same
    if a.Suit == b.Suit && a.RankIndex() == b.RankIndex() {
        return false
    }

    // check trumps first
    if a.Suit == trump && b.Suit == trump {
        if a.RankIndex() >= b.RankIndex() {
            return true
        } else {
            return false
        }
    }

    // check which card followed suit
    if a.Suit == ledSuit && b.Suit != ledSuit {
        return true
    }

    if b.Suit == ledSuit && a.Suit != ledSuit {
        return false
    }

    // suits must be the same
    if a.RankIndex() > b.RankIndex() {
        return true
    } else {
        return false
    }
}

// Simple check to place cards by strength
func CardOrderCompare(a, b Card, trump string) int {
    if trump != "NT" {
        a = ReplaceSuit(a, "TR")
        b = ReplaceSuit(b, "TR")
    }

    if a.Suit == b.Suit && a.RankIndex() == b.RankIndex() {
        return 0
    }

    // check trumps first
    if a.Suit > b.Suit {
        return 1
    } else if b.Suit > a.Suit {
        return -1
    } else {
        // suits must be the same
        if a.RankIndex() > b.RankIndex() {
            return 1
        } else {
            return -1
        }
    }
}

func ReplaceSuit(card Card, trump string) Card {
    if card.Suit == trump {
        return Card{card.Rank, trump}
    }
    return card
}

func NewDeck() []Card {
    deck := []Card{}

    for _, r := range ranks {
        for _, s := range suits {
            deck = append(deck, Card{r, s})
        }
    }

    rand.Shuffle(len(deck), func(i, j int) {
        deck[i], deck[j] = deck[j], deck[i]
    })

    return deck
}
