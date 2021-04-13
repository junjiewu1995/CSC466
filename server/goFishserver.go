package main

import "C"
import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

import . "../CallingUtilities"

type GameConfig struct {
    Gameset struct {
        HostPlayers     string `json:"hostplayers"`
    } `json:"gameset`
}

type GoFishServer struct {
	Mu 				      sync.Mutex
	PlayerCounter 	      int
	dead                  bool
	Players               []Player
	Deck			      []Card
	PlayerTurnIndex       int
	TotalPlayers          int
	WinnerPlayerId        int
	Turn 				  int
	Ready				  bool
	score 				  []int
}

type GameStatusArgs struct {
	MatchPair			  string
}

type GameStatusReply struct {
	CurrentPlayerId       int
	TurnIdx 		      int
	Finished 			  bool
	Winner 				  int
	Players 			  []Player
	Turn 				  int
	Ready				  bool
}


/**
 * Load Configurations
 */
func (gfs *GoFishServer) LoadConfiguration (file string) (GameConfig, error) {
	var config GameConfig
	configFile, err := os.Open(file)
	/* defer to the end and close the config file*/
	defer configFile.Close()
	if err != nil { return config, err }
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}


/**
 * If the player does not have the cards, it would return One Card for the Player
 */
func (gfs *GoFishServer) GoFish () (Card){
	if len(gfs.Deck) == 0 {
		return Card{Value: "-1"}
	}
	card := gfs.Deck[0]
	gfs.Deck = gfs.Deck[1:]
	return card
}


func (gfs *GoFishServer) CallForCards(ask *CardRequest, reply *CardRequestReply) error {
	fmt.Println("Player ask for cards : ", ask.Target)
	gfs.Mu.Lock()
	reply.GoFish = true
	// Loop through target player's hand and find matching cards
	var toRemove []int
	var cardPool = gfs.Players[ask.Target].Hand
	for k, v := range cardPool {
		if v.Value == ask.Value {
			fmt.Println("Find the Match ...")
			reply.GoFish = false
			reply.Cards = append(reply.Cards, v)
			toRemove = append(toRemove, k)
		}
	}
	fmt.Println("I am Here")
	if reply.GoFish { // No card found
		fmt.Println("Not Find Match ...")
		reply.Cards = append(reply.Cards, gfs.GoFish())
	} else { // Target player has 1 or more matching cards
		// remove the target players hand
		sort.Ints(toRemove)
		for i, v := range toRemove {
			gfs.Players[ask.Target].Hand = append(gfs.Players[ask.Target].Hand[:v-i], gfs.Players[ask.Target].Hand[v+1-i:]...)
		}
	}
	gfs.Mu.Unlock()
	fmt.Println("ENd")
	return nil
}


/**
 * The gameOver function would check the Status of the game
 */
func (gfs *GoFishServer) GetStatusOfGame (ask *GameStatusArgs, reply *GameStatusReply) error {
	gfs.Mu.Lock()
	defer gfs.Mu.Unlock()
	fmt.Println("Get Game Status ...")
	reply.CurrentPlayerId  = gfs.PlayerTurnIndex
	reply.Finished = gfs.gameOver()
	reply.Turn = gfs.Turn
	reply.Players = gfs.Players
	reply.Ready = gfs.Ready

	return nil
}

/*
 * Call For the End
*/

func (gfs *GoFishServer) CallForEnd (ask *PlayPairRequest, reply *PlayPairReply) error {

	gfs.Mu.Lock()
	defer gfs.Mu.Unlock()

	// Update the player's matching Pairs / passed arrays of pairs
	if ask.Pair != nil && len(ask.Pair) != 0 {
		gfs.Players[ask.Owner].Pair = append(gfs.Players[ask.Owner].Pair, ask.Pair...)
		gfs.score[ask.Owner] += 1
		fmt.Printf("Player %d my pairs: %v\n\n", ask.Owner,  gfs.Players[ask.Owner].Pair)
	}

	// Update the player's hand
	gfs.Players[ask.Owner].Hand = ask.Hand

	// Determine next player
	gfs.PlayerTurnIndex++
	if gfs.PlayerTurnIndex  >= gfs.TotalPlayers {
		gfs.PlayerTurnIndex  = 0
	}
	ask.Pair = nil
	return nil
}


/**
 * Players Enter Game
*/
func (gfs *GoFishServer) EnterGame (ask *CardRequest, reply *CardRequestReply) error {
    gfs.Mu.Lock()

    if gfs.PlayerCounter < 6 {
        reply.ID = gfs.PlayerCounter
		fmt.Println("Player", gfs.PlayerTurnIndex, "Enter Game ...")
        gfs.Players = append(gfs.Players, Player{ID: gfs.PlayerCounter})
		gfs.score = append(gfs.score, 0)
		fmt.Println("*****")
		for k,v := range(gfs.score){
			fmt.Println(k, v)
		}
		fmt.Println("*****")
        gfs.PlayerCounter += 1
    } else {
		reply.ID = -1
	}

    if gfs.Equals(gfs.PlayerCounter, gfs.TotalPlayers) {
    	/* Loading Cards */
		fmt.Println("Load cards ...")
		time.Sleep(1 * time.Second)
		gfs.LoadCard()
		/* Assign Cards */
		fmt.Println("Assign cards ...")
		time.Sleep(1 * time.Second)
		gfs.assignCard()

		gfs.Ready = true
    }

	gfs.Mu.Unlock()
    return nil
}


/* RPC server interaction */
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


/* Game Over */
func (gfs *GoFishServer) gameOver() bool {

	if len(gfs.Deck) == 0 && gfs.Ready {
		fmt.Println("Game Over ...")
		gfs.dead = true
		max := 0
		for k, v := range gfs.score {
			// winner
			if max < v {
				max = v
				gfs.WinnerPlayerId = k
			}
		}
		fmt.Println(gfs.WinnerPlayerId)
	}

    return gfs.dead
}

/* Create a Game Server */
func StartServer () {
    gfs := GoFishServer{}
    gfs.StateSet()

    /* Construct the Game Config file */
	config, _ := gfs.LoadConfiguration("../game.config.json")
	i1, err := strconv.Atoi(config.Gameset.HostPlayers)
    if err == nil { gfs.TotalPlayers = i1 }

	/* Calling server method */
	gfs.server()

	/* Checking the game is over or not periodly */
    for !gfs.gameOver() { time.Sleep( 1 * time.Second) }
}

func main () {
	StartServer() /* Start to Run Go Fish Server */
}