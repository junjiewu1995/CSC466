package main

import (
    "errors"
    "fmt"
    "math/rand"
    "strings"
    "time"
)

const Spy = false

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Card int

func (c Card) String() string {
    return [...]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}[int(c)]
}

func parseCard(s string) (Card, error) {
    for i := 0; i < 13; i++ {
        card := Card(i)
        if card.String() == s {
            return card, nil
        }
    }
    return Card(0), errors.New("Unknown card value: " + s)
}

type Player struct {
    hand     []int
    score    int
    computer bool
}

func newPlayer(computer bool) *Player {
    return &Player{make([]int, 13, 13), 0, computer}
}

func (p *Player) TakeCard(card Card) {
    p.hand[card]++
    if p.hand[card] == 4 {
        p.hand[card] = 0
        p.score++
        fmt.Printf("There's a complete book of %s.\n", card)
    }
}

func (p *Player) NumberOfCards() int {
    total := 0
    for _, count := range p.hand {
        total += count
    }
    return total
}

type GoFishGame struct {
    deck      []Card
    players   []*Player
    turn      *Player
    turnIndex int
}

func NewGoFishGame(computer ...bool) *GoFishGame {
    deck := make([]Card, 52, 52)
    for i, card := range rnd.Perm(52) {
        deck[i] = Card(card % 13)
    }

    players := make([]*Player, len(computer), len(computer))
    for i, computer := range computer {
        players[i] = newPlayer(computer)
    }
    return &GoFishGame{deck, players, players[0], 0}
}

func (gm *GoFishGame) playerName(index int) string {
    if gm.players[index].computer {
        return fmt.Sprintf("Computer %d", index+1)
    } else {
        return fmt.Sprintf("Human %d", index+1)
    }
}

func (gm *GoFishGame) decideStealingComputer() (int, Card) {
    opponentIndex := rnd.Intn(len(gm.players) - 1)
    if opponentIndex >= gm.turnIndex {
        opponentIndex++
    }

    var choices []Card
    for card, count := range gm.turn.hand {
        for i := 0; i < count; i++ {
            choices = append(choices, Card(card))
        }
    }
    rank := choices[rnd.Intn(len(choices))]

    return opponentIndex, rank
}

func (gm *GoFishGame) decideStealingHuman() (int, Card) {
again:
    var index int
    if len(gm.players) == 2 {
        index = (gm.turnIndex + 1) % 2
    } else {
        var oneIndex int
        fmt.Printf("Steal from whom (1-%d)? ", len(gm.players))
        fmt.Scan(&oneIndex)
        index = oneIndex - 1
        if !(0 <= index && index <= len(gm.players) && index != gm.turnIndex) {
            fmt.Println("Wrong opponent.")
            goto again
        }
    }

    fmt.Print("Steal which rank? ")
    var card string
    fmt.Scan(&card)
    rank, err := parseCard(strings.ToUpper(card))
    if err != nil {
        fmt.Println("Wrong rank.")
        goto again
    }

    if gm.turn.hand[rank] == 0 {
        fmt.Println("You can only steal a rank from your hand.")
        goto again
    }
    return index, rank
}

func (gm *GoFishGame) decideStealing() (int, Card) {
    if gm.turn.computer {
        return gm.decideStealingComputer()
    } else {
        return gm.decideStealingHuman()
    }
}

func (gm *GoFishGame) drawCard(quiet bool) {
    card := gm.deck[0]
    gm.deck = gm.deck[1:]
    gm.turn.TakeCard(card)
    if !quiet {
        if Spy || !gm.turn.computer {
            fmt.Printf("%s drew card %s.\n", gm.playerName(gm.turnIndex), card)
        } else {
            fmt.Printf("%s drew a card.\n", gm.playerName(gm.turnIndex))
        }
    }
}

func (gm *GoFishGame) nextPlayer(quiet bool) {
    gm.turnIndex = (gm.turnIndex + 1) % len(gm.players)

    gm.turn = gm.players[gm.turnIndex]

    if !quiet {
        fmt.Printf("%s to play.\n", gm.playerName(gm.turnIndex))
    }
}

func (gm *GoFishGame) setup() {
    for range gm.players {
        for i := 0; i < 9; i++ {
            gm.drawCard(true)
        }
        gm.nextPlayer(true)
    }
}

func (gm *GoFishGame) gameOver() bool {

    if len(gm.deck) != 0 {
        return false
    }

    for _, player := range gm.players {
        if player.NumberOfCards() > 0 {
            return false
        }
    }

    return true
}

func (gm *GoFishGame) printStatus() {
    for p, player := range gm.players {
        if p != 0 {
            fmt.Print("   ")
        }
        n := player.NumberOfCards()
        fmt.Print(gm.playerName(p), ": ", n, " cards")
        if (Spy || !player.computer) && n > 0 {
            fmt.Print(" (")
            sep := ""
            for rank, count := range player.hand {
                for i := 0; i < count; i++ {
                    fmt.Print(sep, Card(rank))
                    sep = " "
                }
            }
            fmt.Print(")")
        }
        fmt.Print(", ", player.score, " points.")
    }
    fmt.Print("   ", len(gm.deck), " cards left.")
    fmt.Println()
}

func (gm *GoFishGame) move() {
    player := gm.turn

again:
    if player.NumberOfCards() == 0 {
        if len(gm.deck) == 0 {
            return
        }
        gm.drawCard(false)
    }

    opponentIndex, rank := gm.decideStealing()
    if player.hand[rank] == 0 {
        panic("strategy generated illegal move")
    }

    opponent := gm.players[opponentIndex]
    stolen := opponent.hand[rank]
    opponent.hand[rank] = 0

    if stolen != 0 {
        fmt.Printf("%s stole %d cards of rank %s from %s.\n",
            gm.playerName(gm.turnIndex), stolen, rank, gm.playerName(opponentIndex))

        for i := 0; i < stolen; i++ {
            player.TakeCard(rank)
        }
        if player.hand[rank] == 0 && gm.gameOver() {
            return
        }
        gm.printStatus()
        goto again
    }

    fmt.Printf("%s does not have cards of rank %s.\n", gm.playerName(opponentIndex), rank)
    if len(gm.deck) > 0 {
        gm.drawCard(false)
    }
}

func (gm *GoFishGame) Play() {
    gm.setup()

    for !gm.gameOver() {
        gm.printStatus()
        gm.move()
        gm.nextPlayer(false)
    }

    fmt.Println("Game over.")
    gm.printStatus()
}


func main() {
	//5 CPU players
    gm := NewGoFishGame(false, false, false, false, false)
    gm.Play()
}
