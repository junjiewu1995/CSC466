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

	Mu 				      sync.Mutex
	PlayerCounter 	      int
	dead                  bool
	Players               []Player
	Deck			      []Card
	PlayerTurnIndex       int
	TotalPlayers          int

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


/**
 *  Players Set Up
*/

func (gfs *GoFishServer) PlayersSetUp (cardNum int) {
    fmt.Println("2 players setting up begin ...")

    /* Assign 7 Cards for Each Players */
    for range gfs.Players {
        for i := 0; i < cardNum; i++ {
            gfs.drawCards(i)
        }
        /* Increase the TurnIndex Number */
        gfs.PlayerTurnIndex += 1
    }
    gfs.PlayerTurnIndex = 0
}


/*
 * Player Draw Cards
*/

func (gfs *GoFishServer) drawCards(i int) error {

     /* Draw one card from the Deck and renew the Deck */
     card := gfs.Deck[0]
     gfs.Deck = gfs.Deck[1:]
     card.Used = true
     /* Give the Card to the assignPlayer */
     gfs.Players[gfs.PlayerTurnIndex].Hand = append(gfs.Players[gfs.PlayerTurnIndex].Hand, card)

     return nil
}


/**
 * Players Enter Game
*/
func (gfs *GoFishServer) EnterGame (playerask *CardRequest, reply *CardRequestReply) error {

    /* Lock for each players */
    gfs.Mu.Lock()
    defer gfs.Mu.Unlock()

    /* No More than 7 players */
    if gfs.PlayerCounter < 7 {
        /* Add up the Player number */
        reply.ID = gfs.PlayerCounter

        /* Appends the players */
        gfs.Players = append(gfs.Players, Player{ID: gfs.PlayerCounter})

        /* add up the player counter */
        gfs.PlayerCounter += 1
    }
    /* Once the Player meets the decides number Game starts */
    if gfs.PlayerCounter == gfs.TotalPlayers { gfs.gameStart() }
    return nil
}


/*
 * gameStart => LoadCards()
 *           => AssignCards()
*/

func (gfs *GoFishServer) gameStart () {
    fmt.Println("Game Initializes Environments ... ")
    /* Load Cards */
    gfs.LoadCard()
    /* Assign Initialized Cards */
    gfs.assignCard()
}

/**
 * game initilization
*/

func (gfs *GoFishServer) assignCard () error {

    /* Check the Player Number */
    switch {

        case gfs.PlayerCounter == 1:
            gfs.dead = true

        case gfs.PlayerCounter == 2:
            gfs.PlayersSetUp(7) // 2 players assign 7 cards

        default:
            gfs.PlayersSetUp(5) // More than 2 players assign 5 cards
    }
    return nil
}

/**
 * The gameOver function would check the Status of the game
*/

func (gfs *GoFishServer) GetStatusOfGame () {

}


/**
 * If the player does not have the cards, it would return One Card for the Player
 */
func (gfs *GoFishServer) GoFish () {

}


/*
 * Once the Player call RequestForCard for a particular player
 * The server request the tartget Players' information
 * And return back to the Caller or Call Go Fish
*/
func (gfs *GoFishServer) RequestForCard(ask *CardRequest, reply *CardRequestReply) error {

	gfs.Mu.Lock()
	defer gfs.Mu.Unlock()


	reply.GoFishGame = false
	reply.Turn = 1

	return nil
}

//Fills gfs.Deck with 52 shuffled Cards
func (gfs *GoFishServer) LoadCard() error {

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

    /* Construct the Server Struct */
    gfs := GoFishServer{}
    gfs.Deck = []Card{}

    /* Intialize the Game Variables */
    gfs.TotalPlayers = 0
    gfs.PlayerCounter = 0
    gfs.PlayerTurnIndex = 0

	return gfs
}

/**
 * gameOver
 */

func (gfs *GoFishServer)gameOver() bool {

    return gfs.dead
}

/* Create a Game Server */
func StartServer () *GoFishServer {

    /* Construct the Game Config file */
	config, _ := gfs.LoadConfiguration("../game.config.json")

	i1, err := strconv.Atoi(config.Gameset.HostPlayers)
    if err == nil { gfs.TotalPlayers = i1 }

	/* Calling server method */
	gfs.server()
    rep := gfs.serverStateSet()

    /* Checking the game is over or not periodly */
    for !gfs.gameOver() {
        fmt.Println("Check the Game Status in 1 sec ...")
        /* Check the game status */
        gfs.GetStatusOfGame()
        time.Sleep( 1 * time.Second)
    }

	return rep
}

/* Main Function */
func main () {
    /* Start to Run Go Fish Server */
	StartServer()
}