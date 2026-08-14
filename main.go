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

type DrH struct {
    Hand []Card
}

type East struct {
    Hand []Card
}

type West struct {
    Hand []Card
}

type Teddy struct {
    Hand []Card
}

// Player hands
func (d *DrH) GetHand() []Card {
    return d.Hand
}

func (e *East) GetHand() []Card {
    return e.Hand
}

func (w *West) GetHand() []Card {
    return w.Hand
}

func (t *Teddy) GetHand() []Card {
    return t.Hand
}

// Player Labels
func (d *DrH) GetLabel() string {
    return "Dr.H."
}

func (e *East) GetLabel() string {
    return "East."
}

func (w *West) GetLabel() string {
    return "West."
}

func (t *Teddy) GetLabel() string {
    return "Teddy"
}

// Interfaces
type Player interface {
    GetHand() []Card
    GetLabel() string
}

type Bid struct {
    Tricks int
    Suit string
}

type Game struct {
    Declarer Player
    Board Player

    DrHScore int
    EastWestScore int
}

func (g Game) GetDeclarer() Player {
    return g.Declarer
}

func main() {
    fmt.Println("Dr. H. Bridge")

    deck := []Card{
        Card{2, "C"}, Card{3, "C"}, Card{14, "S"},
    }

    drh := DrH{[]Card{Card{8, "C"}}}
    east := East{[]Card{}}
    west := West{[]Card{}}
    teddy := Teddy{[]Card{}}

    game := Game{
        Declarer: nil,
        Board: nil,

        DrHScore: 0,
        EastWestScore: 0,
    }

    if game.Declarer != nil {
        fmt.Printf("Deck: %v\n", deck)

        fmt.Printf("DrH: %v\n", drh)
        fmt.Printf("East: %v\n", east)
        fmt.Printf("West: %v\n", west)
        fmt.Printf("Teddy: %v\n", teddy)

        fmt.Printf("Game.Declarer: %v\n", game.Declarer)
    }
    // set Declarer
    game.Declarer = &drh
    fmt.Printf("Game.Declarer after reassignment: %v\n", game.Declarer)

    drh.Hand = append(drh.Hand, Card{12, "C"})

    fmt.Printf("Game.Declarer after append: %v\n", game.Declarer)
    fmt.Printf("Game.Declarer from GetDeclarer: %v\n", game.GetDeclarer())
    fmt.Printf("Game.Declarer.Hand from GetDeclarer: %v\n", game.GetDeclarer().GetHand())
    fmt.Printf("Game.Declarer.Label: %v\n", game.GetDeclarer().GetLabel())

    fmt.Println("end of script")
}
