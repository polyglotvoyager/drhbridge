package main

import (
    "fmt"
)

type DrH struct {
    Hand []Card
    Client *Client
}

type West struct {
    Hand []Card
    Client *Client
}

type Teddy struct {
    Hand []Card
    Client *Client
}

type East struct {
    Hand []Card
    Client *Client
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

// Get Client
func (d *DrH) GetClient() *Client {
    return d.Client
}

func (w *West) GetClient() *Client {
    return w.Client
}

func (e *East) GetClient() *Client {
    return e.Client
}

func (t *Teddy) GetClient() *Client {
    return t.Client
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

// Set Client
func (d *DrH) SetClient(newClient *Client) {
    fmt.Println("setting DrH client")
    if d.Client != nil {
        d.Client.hub.unregister <- d.Client
        d.Client.conn.Close()
    }
    d.Client = newClient
}

func (w *West) SetClient(newClient *Client) {
    fmt.Println("setting West client")
    w.Client = newClient
}

func (e *East) SetClient(newClient *Client) {
    fmt.Println("setting East client")
    e.Client = newClient
}

func (t *Teddy) SetClient(newClient *Client) {
    fmt.Println("setting Teddy client")
    t.Client = newClient
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
    GetClient() *Client
    SetClient(*Client)
}
