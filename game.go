package main

import (
    "fmt"
    "strings"
    "strconv"
    "slices"
)

type Bid struct {
    Tricks int
    Suit string
}

// lead is repeated, used to obtain the lead suit
// ignore the possibility that a player cheated (played a trump when he or she had the led suit)
func (g Game) SetTrickWinner() {
    leadSuit := g.TrickLeadCard.Suit
    trump := g.LastBid.Suit

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

func (g Game) PlayerHand(playerName string) []Card {
    switch playerName {
    case "teddy":
        return g.Teddy.GetHand()
    case "drh":
        return g.DrH.GetHand()
    case "west":
        return g.West.GetHand()
    default:
        return nil
    }
}

type Game struct {
    West Player
    East Player
    DrH Player
    Teddy Player

    WestActive bool
    DrHActive bool
    EastActive bool

    Declarer Player
    Dummy Player

    DrHScore int
    EastWestScore int

    Bidder Player
    LastBid Bid

    TrickLeadCard Card
    TrickWinner Player
    TrickWinningCard Card
    TrickDrH Card
    TrickWest Card
    TrickTeddy Card
    TrickEast Card

    TotalTricksDrH int
    TotalTricksEastWest int

    GameState string // keep track if "bidding" or "playing"
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

func (g Game) Play(playerName string, cardString string) string {
    var player Player

    switch playerName {
    case "teddy":
        player = g.Teddy
    case "drh":
        player = g.DrH
    case "west":
        player = g.West
    case "east":
        player = g.East
    }

    cardParts := strings.Split(cardString, " ")
    card := Card{cardParts[0], cardParts[1]}

    index := slices.Index(player.GetHand(), card)

    if index >= 0 {
        player.SetHand(slices.Delete(player.GetHand(), index, index + 1))
        return playerName + " played " + cardString
    } else {
        return playerName + ": card " + cardString + " not found"
    }
}

func (g *Game) Reset() {
    deck := NewDeck()

    g.DrH.SetHand(deck[0:13])
    g.West.SetHand(deck[13:26])
    g.Teddy.SetHand(deck[26:39])
    g.East.SetHand(deck[39:52])

    // sort
    slices.SortFunc(g.DrH.GetHand(), func(a, b Card) int {
        return CardOrderCompare(a, b, "NT") // No Trump for initial sort
    })
    slices.SortFunc(g.West.GetHand(), func(a, b Card) int {
        return CardOrderCompare(a, b, "NT") // No Trump for initial sort
    })
    slices.SortFunc(g.Teddy.GetHand(), func(a, b Card) int {
        return CardOrderCompare(a, b, "NT") // No Trump for initial sort
    })
    slices.SortFunc(g.East.GetHand(), func(a, b Card) int {
        return CardOrderCompare(a, b, "NT") // No Trump for initial sort
    })

    g.Bidder = g.DrH // DrH always goes first in a newly created game
    g.GameState = "bidding"
}

func (g Game) NextPlayer(player Player) Player {
    switch player.GetLabel() {
    case "Dr.H.":
        return g.West
    case "West":
        return g.Teddy
    case "Teddy":
        return g.East
    case "East":
        return g.DrH
    default:
        return nil
    }
}

func NewGame() Game {
    g := Game{}

    drh := &DrH{Hand: make([]Card, 13)}
    west := &West{Hand: make([]Card, 13)}
    teddy := &Teddy{Hand: make([]Card, 13)}
    east := &East{Hand: make([]Card, 13)}

    g.DrH = drh
    g.West = west
    g.Teddy = teddy
    g.East = east

    g.Bidder = g.DrH

    return g
}

func (g *Game) PlayerActive(player string) {
    switch player {
    case "west":
        g.WestActive = true
    case "east":
        g.EastActive = true
    case "drh":
        g.DrHActive = true
    }
}

func (g *Game) PlayerInactive(player string) {
    switch player {
    case "west":
        g.WestActive = false
    case "east":
        g.EastActive = false
    case "drh":
        g.DrHActive = false
    }
}

func (g Game) GetPlayer(name string) Player {
    switch name {
    case "west":
        return g.West
    case "east":
        return g.East
    case "drh":
        return g.DrH
    case "teddy":
        return g.Teddy
    default:
        return nil
    }
}

func (g Game) AllActive() bool {
    return g.WestActive && g.EastActive && g.DrHActive
}

func (g *Game) Bid(playerName string, bidString string) string {
    player := g.GetPlayer(playerName)
    bidParts := strings.Split(bidString, " ")
    tricks, err := strconv.Atoi(bidParts[0])
    if err != nil {
        return "Invalid number of tricks: " + bidParts[0]
    }
    suit := bidParts[1]
    b := Bid{tricks, suit}

    g.LastBid = b
    player.SetBid(b)
    g.Bidder = g.NextPlayer(player)

    return fmt.Sprintf(player.GetLabel() + " bid %v", g.LastBid)
}

func (g Game) SetPlayerClient(username string, client *Client) {
    playerName := strings.ReplaceAll(username, "debug", "")
    switch playerName {
    case "drh":
        g.DrH.SetClient(client)
    case "west":
        g.West.SetClient(client)
    case "east":
        g.East.SetClient(client)
    case "teddy":
        g.Teddy.SetClient(client)
    default:
        fmt.Println("SetPlayerClient: unknown player " + playerName)
    }
}
