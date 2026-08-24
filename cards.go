package main

import (
    "slices"
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

func ReplaceSuit(card Card, trump string) Card {
    if card.Suit == trump {
        return Card{card.Rank, trump}
    }
    return card
}
