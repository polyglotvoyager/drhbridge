// to test locally: start username with "debug"

package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "math/rand/v2"
    "os"
    "embed"
)

//go:embed templates
var templates embed.FS

var addr = flag.String("addr", ":8100", "http service address")

func serveRegister(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, templates, "templates/register.html")
}

func servePlay(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, templates, "templates/play.html")
}

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
    West []Card
    East []Card
    DrH []Card
    Teddy []Card

    Declarer Player
    Board Player

    DrHScore int
    EastWestScore int
}

func (g Game) GetDeclarer() Player {
    return g.Declarer
}

func (g Game) WestHand() []Card {
    return g.West
}

func (g Game) EastHand() []Card {
    return g.East
}

func (g Game) DrHHand() []Card {
    return g.DrH
}

func (g Game) TeddyHand() []Card {
    return g.Teddy
}

func NewGame() Game {
    deck := NewDeck()

    return Game{
        West: deck[0:13],
        East: deck[13:26],
        DrH: deck[26:39],
        Teddy: deck[39:52],

        Declarer: nil,
        Board: nil,

        DrHScore: 0,
        EastWestScore: 0,
    }
}

func NewDeck() []Card {
    deck := []Card{}
    suits := []string{"C", "D", "H", "S"}

    for r := 2; r <= 14; r++ {
        for sIdx := 0; sIdx < 4; sIdx++ {
            deck = append(deck, Card{r, suits[sIdx]})
        }
    }

    rand.Shuffle(len(deck), func(i, j int) {
        deck[i], deck[j] = deck[j], deck[i]
    })

    return deck
}

func main() {
    fmt.Println("Dr. H. Bridge")

    game := NewGame()

    flag.Parse()
    hub := newHub(game)
    go hub.run()

    mux := http.NewServeMux()
    mux.HandleFunc("GET /drhbridge", serveRegister)
    mux.HandleFunc("GET /drhbridge/{$}", serveRegister)
    mux.HandleFunc("GET /drhbridge/play", servePlay)
    mux.HandleFunc("/drhbridge/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    srv := &http.Server{
        Addr: *addr,
        Handler: mux,
    }
    fmt.Printf("serving drhbridge at %v\n", *addr)
    err := srv.ListenAndServe()
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
    os.Exit(1)
}
