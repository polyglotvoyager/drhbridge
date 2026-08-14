package main

import (
    "fmt"
)

// rank: 2 to 14 (J=11, Q=12, K=13, A=14)
// suit: C, D, H, S
type Card struct {
    Rank int
    Suit string
}

type Player struct {
    Hand []Card
    IsDrH bool
}

type Bid struct {
    Tricks int
    Suit string
}

type Game struct {
    Declarer *Player
    Board *Player
    EastWestScore int
    DrHScore int
}

func (c Card) Write() {

}

func main() {
    fmt.Println("Dr. H. Bridge")

    deck := []Card{
        Card{2, "C"}, Card{3, "C"}, Card{14, "S"},
    }

    fmt.Printf("Deck: %v\n", deck)

    fmt.Println("end of script")
}
