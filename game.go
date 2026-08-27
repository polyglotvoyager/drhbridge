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
    Special string
}

func (b Bid) GetLabel() string {
    return strconv.Itoa(b.Tricks) + b.Suit
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
    case "east":
        return g.East.GetHand()
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

    TrickPlayer Player
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

    hub *Hub
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

func (g *Game) Play(playerName string, cardString string) string {
    player := g.GetPlayer(playerName)

    if player != g.TrickPlayer {
        return player.GetLabel() + " played out of turn"
    }

    cardParts := strings.Split(cardString, " ")
    card := Card{cardParts[0], cardParts[1]}

    index := slices.Index(player.GetHand(), card)

    if index >= 0 {
        player.SetHand(slices.Delete(player.GetHand(), index, index + 1))
        g.NextPlayer()
        // TODO: handle last trick (count TotalTricksDrH and TotalTricksEastWest), if == 13, stop
        g.GameCommand()
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

func (g *Game) NextBidder() {
    switch g.Bidder {
    case g.DrH:
        g.Bidder = g.West
    case g.West:
        g.Bidder = g.Teddy
    case g.Teddy:
        g.Bidder = g.East
    case g.East:
        g.Bidder = g.DrH
    default:
        panic("could not advance bidder")
    }
}

func (g *Game) NextPlayer() {
    switch g.TrickPlayer {
    case g.DrH:
        g.TrickPlayer = g.West
    case g.West:
        g.TrickPlayer = g.Teddy
    case g.Teddy:
        g.TrickPlayer = g.East
    case g.East:
        g.TrickPlayer = g.DrH
    default:
        panic("could not advance player")
    }
}

func NewGame() *Game {
    g := &Game{}

    drh := &DrH{Hand: make([]Card, 13)}
    west := &West{Hand: make([]Card, 13)}
    teddy := &Teddy{Hand: make([]Card, 13)}
    east := &East{Hand: make([]Card, 13)}

    g.DrH = drh
    g.West = west
    g.Teddy = teddy
    g.East = east

    g.Bidder = g.DrH
    g.LastBid = Bid{0, "0", "new"}

    g.hub = newHub()

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


// if the player is the third pass in a row, return true
func (g *Game) AllOthersPass(player Player) bool {
    if player == g.DrH {
        return g.DrH.GetBid().Special == "pass" &&
            g.East.GetBid().Special == "pass" &&
            g.Teddy.GetBid().Special == "pass"
    }
    if player == g.West {
        return g.West.GetBid().Special == "pass" &&
            g.DrH.GetBid().Special == "pass" &&
            g.East.GetBid().Special == "pass"
    }
    if player == g.East {
        return g.East.GetBid().Special == "pass" &&
            g.West.GetBid().Special == "pass" &&
            g.Teddy.GetBid().Special == "pass"
    }
    if player == g.Teddy {
        return g.Teddy.GetBid().Special == "pass" &&
            g.DrH.GetBid().Special == "pass" &&
            g.West.GetBid().Special == "pass"
    }
    return false
}

func (g *Game) AllPass() bool {
    return g.DrH.GetBid().Special == "pass" &&
        g.West.GetBid().Special == "pass" &&
        g.East.GetBid().Special == "pass" &&
        g.Teddy.GetBid().Special == "pass"
}

func (g *Game) Bid(playerName string, bidString string) string {
    player := g.GetPlayer(playerName)

    if player != g.Bidder {
        return player.GetLabel() + " played out of turn"
    }

    bidParts := strings.Split(bidString, " ")
    tricks, err := strconv.Atoi(bidParts[0])
    if err != nil {
        return "Invalid number of tricks: " + bidParts[0]
    }
    suit := bidParts[1]
    special := ""

    if len(bidParts) == 3 {
        special = bidParts[2]
    }

    b := Bid{tricks, suit, special}

    g.LastBid = b
    player.SetBid(b)

    g.NextBidder()

    if g.AllPass() {
        return "all pass - TODO - reshuffle"
    }

    if g.AllOthersPass(player) {
        g.Declarer = g.Bidder
        g.LastBid = g.Declarer.GetBid()
        g.TrickPlayer = g.Declarer
        g.GameState = "playing"
    }
    g.GameCommand()
    return fmt.Sprintf(player.GetLabel() + " bid %v", player.GetBid())
}

func (g *Game) SetPlayerClient(username string, client *Client) {
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

func (g Game) GameCommand() {
    message := "GAME: " + g.GameState + " "
    if g.GameState == "bidding" {
        message += g.Bidder.GetLabel()
    } else if g.GameState == "playing" {
        message += g.LastBid.GetLabel() + " " + g.TrickPlayer.GetLabel() + " to play"
    } else {
        message += "(invalid state)"
    }
    g.hub.broadcast <- []byte(message)
}
