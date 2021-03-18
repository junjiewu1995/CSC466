package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"math/rand"
	"time"
)

import . "../helper"


type Player struct {
	ID 			int
	Hand 		[]Card
	Pairs 		[]Pairs
	Opponents 	[]Player
}

type GoFishServer struct {
	Deck			[]Card
	Mu 				sync.Mutex
	PlayerCounter 	int
	dead            bool
}

func (gfs *GoFishServer) RequestForCard(ask *CardRequest, reply *CardRequestReply) error {

	gfs.Mu.Lock()
	defer gfs.Mu.Unlock()

    fmt.Println("Calling from RequestForCard ...")
	reply.GoFishGame = false
	reply.Turn = 1

	return nil
}

//Fills gfs.Deck with 52 shuffled Cards
func (gfs *GoFishServer) LoadCard() {
	//ensure Deck is empty to start
	gfs.Deck := []Card{}
	
	//values a card can be
	cardValues := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	//Create 52 new cards, not shuffled, in gfs.Deck
	for i := 0; i < 13; i++ {
		gfs.Deck = append(gfs.Deck, Card{Value: cardValues[i], Suit: "clubs"})
	}
	for i := 0; i < 13; i++ {
		gfs.Deck = append(gfs.Deck, Card{Value: cardValues[i], Suit: "diamonds"})
	}
	for i := 0; i < 13; i++ {
		gfs.Deck = append(gfs.Deck, Card{Value: cardValues[i], Suit: "hearts"})
	}
	for i := 0; i < 13; i++ {
		gfs.Deck = append(gfs.Deck, Card{Value: cardValues[i], Suit: "spades"})
	}

	//shuffle gfs.Deck
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(gfs.Deck), func(i, j int) { gfs.Deck[i], gfs.Deck[j] = gfs.Deck[j], gfs.Deck[i] })
}

func (gfs *GoFishServer) printLine() {
	fmt.Println("Hello, World :)")
}

/**
 * RPC server interaction
*/
func (gfs *GoFishServer) server() {
	rpc.Register(gfs)
	rpc.HandleHTTP()
	sockname := MasterSock()
	os.Remove(sockname)
    l, e := net.Listen("unix", sockname)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    go http.Serve(l, nil)
}


func (gfs *GoFishServer)serverStateSet() *GoFishServer {
	return gfs
}



/* Create a Game Server */
func StartServer () *GoFishServer {

    /* Construct the Server Struct */
	gfs := GoFishServer{}

	/* Calling server method */
	gfs.server()

    rep := gfs.serverStateSet()

    /* Checking the game is over or not */
    for gfs.dead == false {

    }

	return rep
}

/* Main Function */
func main () {
    /* Start to Run Go Fish Server */
	StartServer()
}