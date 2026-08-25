package main

type DrH struct {
    Hand []Card
}

type West struct {
    Hand []Card
}

type Teddy struct {
    Hand []Card
}

type East struct {
    Hand []Card
}

// Player hands
func (d *DrH) GetHand() []Card {
    return d.Hand
}

func (w *West) GetHand() []Card {
    return w.Hand
}

func (t *Teddy) GetHand() []Card {
    return t.Hand
}

func (e *East) GetHand() []Card {
    return e.Hand
}

// Setters
func (d *DrH) SetHand(newHand []Card) {
    d.Hand = newHand
}

func (w *West) SetHand(newHand []Card) {
    w.Hand = newHand
}

func (t *Teddy) SetHand(newHand []Card) {
    t.Hand = newHand
}

func (e *East) SetHand(newHand []Card) {
    e.Hand = newHand
}

// Player Labels
func (d *DrH) GetLabel() string {
    return "Dr.H."
}

func (w *West) GetLabel() string {
    return "West."
}

func (t *Teddy) GetLabel() string {
    return "Teddy"
}

func (e *East) GetLabel() string {
    return "East."
}

// Interfaces
type Player interface {
    GetHand() []Card
    SetHand([]Card)
    GetLabel() string
}
