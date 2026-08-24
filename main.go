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
func (d DrH) GetHand() []Card {
    return d.Hand
}

func (w West) GetHand() []Card {
    return w.Hand
}

func (t Teddy) GetHand() []Card {
    return t.Hand
}

func (e East) GetHand() []Card {
    return e.Hand
}

// Player Labels
func (d DrH) GetLabel() string {
    return "Dr.H."
}

func (w West) GetLabel() string {
    return "West."
}

func (t Teddy) GetLabel() string {
    return "Teddy"
}

func (e East) GetLabel() string {
    return "East."
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

// lead is repeated, used to obtain the lead suit
// ignore the possibility that a player cheated (played a trump when he or she had the led suit)
func (g Game) SetTrickWinner() {
    leadSuit := g.TrickLeadCard.Suit
    trump := g.WinningBid.Suit

    // replace cards' suits with trump if it exists
    if trump != "NT" {
        g.TrickDrH = ReplaceSuit(g.TrickDrH, trump)
        g.TrickWest = ReplaceSuit(g.TrickWest, trump)
        g.TrickTeddy = ReplaceSuit(g.TrickTeddy, trump)
        g.TrickEast = ReplaceSuit(g.TrickEast, trump)
    }

    if CardBeats(g.TrickDrH, g.TrickLeadCard, leadSuit, trump) {
        g.TrickWinner = g.DrH
        g.TrickWinningCard = g.TrickDrH
    }

    if CardBeats(g.TrickWest, g.TrickWinningCard, leadSuit, trump) {
        g.TrickWinner = g.West
        g.TrickWinningCard = g.TrickWest
    }

    if CardBeats(g.TrickTeddy, g.TrickWinningCard, leadSuit, trump) {
        g.TrickWinner = g.Teddy
        g.TrickWinningCard = g.TrickTeddy
    }

    if CardBeats(g.TrickEast, g.TrickWinningCard, leadSuit, trump) {
        g.TrickWinner = g.East
        g.TrickWinningCard = g.TrickEast
    }
}

type Game struct {
    West Player
    East Player
    DrH Player
    Teddy Player

    Declarer Player
    Dummy Player

    DrHScore int
    EastWestScore int

    WinningBid Bid

    TrickLeadCard Card
    TrickWinner Player
    TrickWinningCard Card
    TrickDrH Card
    TrickWest Card
    TrickTeddy Card
    TrickEast Card
}

func (g Game) GetDeclarer() Player {
    return g.Declarer
}

func (g Game) DrHHand() []Card {
    return g.DrH.GetHand()
}

func (g Game) WestHand() []Card {
    return g.West.GetHand()
}

func (g Game) TeddyHand() []Card {
    return g.Teddy.GetHand()
}

func (g Game) EastHand() []Card {
    return g.East.GetHand()
}

func NewGame() Game {
    deck := NewDeck()

    return Game{
        West: West{deck[0:13]},
        East: East{deck[13:26]},
        DrH: DrH{deck[26:39]},
        Teddy: Teddy{deck[39:52]},

        Declarer: nil,
        Dummy: nil,

        DrHScore: 0,
        EastWestScore: 0,
    }
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

    fmt.Printf("drh hand %v\n", game.DrHHand())
    err := srv.ListenAndServe()
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
    os.Exit(1)
}
