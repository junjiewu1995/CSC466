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
	"encoding/json"
	"strconv"
)

import . "../CallingUtilities"


type GameConfig struct {
    Gameset struct {
        HostPlayers     string `json:"hostplayers"`
    } `json:"gameset`
}


type Player struct {
	ID 			int
	Hand 		[]Card
	Pairs 		[]Pairs
	Opponents 	[]Player
}

type GoFishServer struct {

	Mu 				sync.Mutex
	PlayerCounter 	int
	dead            bool
	Players         []Player
	Deck			[]Card
	turnIndex       int
	TotalPlayers    int
}


func (gfs *GoFishServer)LoadConfiguration (file string) (GameConfig, error) {
    var config GameConfig
    configFile, err := os.Open(file)
    /* defer to the end and close the config file*/
    defer configFile.Close()
    if err != nil {
        return config, err
    }
    jsonParser := json.NewDecoder(configFile)
    err = jsonParser.Decode(&config)
    return config, err
}

func (gfs *GoFishServer) EnterGame (playerask *CardRequest, reply *CardRequestReply) error {

    /* Lock for each players */
    gfs.Mu.Lock()
    defer gfs.Mu.Unlock()

    var p Player

    /* No More than 7 players */
    if gfs.PlayerCounter < 7 {

        /* assign player info */
        p.ID = gfs.PlayerCounter

        /* Add up the Player number */
        reply.ID = gfs.PlayerCounter

        /* Appends the players */
        gfs.Players = append(gfs.Players, p)

        /* add up the player counter */
        gfs.PlayerCounter += 1
    }

    /* Once the Player meet the decides number Game starts */
    if gfs.PlayerCounter == gfs.TotalPlayers {
        fmt.Println("Game Starts ... ")
        gfs.gameStart()
    }

    return nil
}

/*
 * gameStart => LoadCards()
 *           => AssignCards()
*/

func (gfs *GoFishServer) gameStart () {
    gfs.LoadCard()
    gfs.assignCard()
}

func (gfs *GoFishServer) assignCard () error {

    for idx, _ := range gfs.Players {
        for x, _ := range gfs.Deck {
           var num = 0
           if num < 5 {
              gfs.Players[idx].Hand = append(gfs.Players[idx].Hand, CardValue)
              num = num + 1
           }
        }
    }
    return nil
}

func (gfs *GoFishServer) GetStatusOfGame () {

}

func (gfs *GoFishServer) RequestForCard(ask *CardRequest, reply *CardRequestReply) error {

	gfs.Mu.Lock()
	defer gfs.Mu.Unlock()

	reply.GoFishGame = false
	reply.Turn = 1

	return nil
}

//Fills gfs.Deck with 52 shuffled Cards
func (gfs *GoFishServer) LoadCard() error {

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
	rand.Shuffle(len(gfs.Deck), func(i, j int) {gfs.Deck[i], gfs.Deck[j] = gfs.Deck[j], gfs.Deck[i] })
    return nil
}

//Draw a card from the deck (returns card type)
func (gfs *GoFishServer) goFish() {
    
    //draw a card from the deck, removing it from gfs.Deck
    drawnCard = gfs.Deck[0]
    gfs.Deck = gfs.Deck[1:]
    
    return drawnCard
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
	gfs.Deck = []Card{}

    gfs.TotalPlayers = 0
	gfs.PlayerCounter = 0

	config, _ := gfs.LoadConfiguration("../game.config.json")

	i1, err := strconv.Atoi(config.Gameset.HostPlayers)
    if err == nil { gfs.TotalPlayers = i1 }

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