package main

import (
	. "../CallingUtilities"
	"math/rand"
	"time"
)


func (gfs *GoFishServer) Equals(pc int, tc int) bool { return pc == tc }

/**
 * Game initilization
 */
func (gfs *GoFishServer) assignCard () error {

	/* Check the Player Number */
	switch {
		case gfs.PlayerCounter == 1:
			gfs.dead = true
		case gfs.PlayerCounter == 2:
			gfs.PlayersSetUp(7)
		default:
			gfs.PlayersSetUp(5)
	}
	return nil
}

/*
 * Player Draw Cards
 */
func (gfs *GoFishServer) drawCards(i bool) error {

	/* Draw one card from the Deck and renew the Deck */
	card := gfs.Deck[0]
	gfs.Deck = gfs.Deck[1:]
	gfs.Players[gfs.PlayerTurnIndex].Hand = append(gfs.Players[gfs.PlayerTurnIndex].Hand, card)

	return nil
}


/*
 * Players Set Up
*/
func (gfs *GoFishServer) PlayersSetUp (cardNum int) {
	for range gfs.Players {
		for i := 0; i < cardNum; i++ {
			gfs.drawCards(true)
		}
		gfs.PlayerTurnIndex += 1
	}
	gfs.PlayerTurnIndex = 0
}


/**
 * LoadCard Methods
*/
func (gfs *GoFishServer) LoadCard() error {
	cardValues := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	//Create 52 new cards, not shuffled, in gfs.Deck
	for i := 0; i < 13; i++ {
		gfs.Deck = append(gfs.Deck, Card{Value: cardValues[i], Suit: "clubs"})
		gfs.Deck = append(gfs.Deck, Card{Value: cardValues[i], Suit: "diamonds"})
		gfs.Deck = append(gfs.Deck, Card{Value: cardValues[i], Suit: "hearts"})
		gfs.Deck = append(gfs.Deck, Card{Value: cardValues[i], Suit: "spades"})
	}
	//shuffle gfs.Deck
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(gfs.Deck), func(i, j int) {gfs.Deck[i], gfs.Deck[j] = gfs.Deck[j], gfs.Deck[i] })
	return nil
}

/**
 * State Set
*/
func (gfs *GoFishServer)StateSet () *GoFishServer {
	gfs.Deck = []Card{}
	gfs.TotalPlayers = 0
	gfs.PlayerCounter = 0
	gfs.PlayerTurnIndex = 0
	gfs.Ready = false // game starts
	return gfs
}

//func (gfs *GoFishServer) saveGameState() {
//	args := GameStateArgs{}
//	reply := GameStateReply{}
//	args.Key = string(gfs.ServerId)
//	js, _ := json.Marshal(gfs)
//	args.Payload = string(js)
//	ok := CallRB("RaftBroker.PutGameState", &args, &reply)
//	if !ok || !reply.Ok {
//		fmt.Printf("Put Game state failed\n")
//	}
//}
//
//func (gs *GoFishServer) getGameState() {
//	gs.Mu.Lock()
//	defer gs.Mu.Unlock()
//	args := GameStateArgs{}
//	reply := GameStateReply{}
//	args.Key = string(gfs.ServerId)
//	ok := CallRB("RaftBroker.GetGameState", &args, &reply)
//	if !ok || !reply.Ok {
//		fmt.Printf("Get Game state failed\n")
//	}
//	gfs.reconcileState(reply.Payload)
//}
//
//
//func (gs *GoFishServer) reconcileState(payload string) {
//	var gsSaved GameServer
//	err := json.Unmarshal([]byte(payload), &gsSaved)
//	if err != nil {
//		fmt.Printf("Unmarshall of game state failed")
//	}
//	if gsSaved.ServerId != gs.ServerId {
//		fmt.Printf("Wrong game state retreived")
//	} else {
//		gs.Winner = gsSaved.Winner
//		gs.Players = gsSaved.Players
//		gs.GameOver = gsSaved.GameOver
//		gs.CurrentTurnPlayer = gsSaved.CurrentTurnPlayer
//		gs.CurrentTurn = gsSaved.CurrentTurn
//		gs.Deck = gsSaved.Deck
//		gs.Ready = gsSaved.Ready
//		gs.GameInitialized = gsSaved.GameInitialized
//	}
//
//}