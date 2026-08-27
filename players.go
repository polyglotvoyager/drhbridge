package main

import (
    "fmt"
    "sync"
)

type DrH struct {
    Hand []Card
    Bid Bid
    Client *Client
    mu sync.Mutex
}

type West struct {
    Hand []Card
    Bid Bid
    Client *Client
    mu sync.Mutex
}

type Teddy struct {
    Hand []Card
    Bid Bid
    Client *Client
    mu sync.Mutex
}

type East struct {
    Hand []Card
    Bid Bid
    Client *Client
    mu sync.Mutex
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

// Bids
func (d *DrH) GetBid() Bid {
    return d.Bid
}

func (t *Teddy) GetBid() Bid {
    return t.Bid
}

func (e *East) GetBid() Bid {
    return e.Bid
}

func (w *West) GetBid() Bid {
    return w.Bid
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

// Bids
func (d *DrH) SetBid(newBid Bid) {
    d.Bid = newBid
}

func (e *East) SetBid(newBid Bid) {
    e.Bid = newBid
}

func (w *West) SetBid(newBid Bid) {
    w.Bid = newBid
}

func (t *Teddy) SetBid(newBid Bid) {
    t.Bid = newBid
}

// Set Client
func (d *DrH) SetClient(newClient *Client) {
    fmt.Println("setting DrH client")

    d.mu.Lock()
    if d.Client != nil {
        d.Client.hub.unregister <- d.Client
        d.Client.conn.Close()
    }
    d.Client = newClient
    d.mu.Unlock()
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
    return "Dr.H"
}

func (w *West) GetLabel() string {
    return "West"
}

func (t *Teddy) GetLabel() string {
    return "Teddy"
}

func (e *East) GetLabel() string {
    return "East"
}

// Interfaces
type Player interface {
    GetHand() []Card
    SetHand([]Card)
    GetLabel() string
    GetClient() *Client
    SetClient(*Client)
    GetBid() Bid
    SetBid(Bid)
}
