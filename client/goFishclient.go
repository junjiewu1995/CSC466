package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	_ "strconv"
	"time"
)

import . "../CallingUtilities"

type Player struct {
	ID 					 int
	Hand 				 []Card
	Pair 				 []Pairs
	Players  			 []Player
}

type GameStatusReply struct {
	CurrentPlayerId       int
	TurnId   		      int
	Finished 			  bool
	Winner 				  int
	Players 			  []Player
	Turn 				  int
	Ready				  bool
}

type GameStatusargs struct {
	MatchPair			  string
}


//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	sockname := MasterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil { log.Fatal("dialing:", err) }
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil { return true }
	fmt.Println(err)
	return false
}


func (p *Player) callForEnd(pairs []Pairs) error {
	reply := CardRequestReply{}
	args := PlayPairRequest{ Owner: p.ID, Pair: pairs, Hand: p.Hand }
	if !call("GoFishServer.CallForEnd", args, &reply) {
		fmt.Println("Fail to request for Cards")
	}
	return nil
}


/**
 * Players Hands Updates
*/
func (p *Player) HandsUpdates () ([]Pairs, error) {
	time.Sleep(1 * time.Second)
	var pairList [] Pairs

	fmt.Println("***", p.Hand, "***")

	for i := 0; i < len(p.Hand); i++ {
		for x := i; x < len(p.Hand); x++ {
			if p.Hand[i].Value == p.Hand[x].Value && !p.Hand[x].Used && !p.Hand[i].Used && i != x {
				p.Hand[i].Used = true
				p.Hand[x].Used = true
				pairList = append(pairList, Pairs{One: p.Hand[i], Two: p.Hand[x]})
			}
		}
	}

	// remove pairs from player's hand
	var newList []Card
	for _, v := range p.Hand {
		if v.Used != true {
			newList = append(newList, v)
		}
	}

	p.Hand = newList
	time.Sleep(3 * time.Second)
    return pairList, nil
}

/*
 * Random Index
*/
func (p *Player) callForStart() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// randomly choose an opponent to ask
	var randIdx = p.ID
	if len(p.Players) > 1 {
		fmt.Printf("Player %d my hand: %v\n", p.ID, p.Hand)
		fmt.Printf("Player %d my pairs: %v\n\n", p.ID, p.Pair)
		// choose a random players
		for randIdx == p.ID { randIdx = r.Intn(len(p.Players)) }
		args := CardRequest{Target: p.Players[randIdx].ID}
		// randomly choose a card value from your hand to ask for
		if len(p.Hand) > 0 {
			randIdx = r.Intn(len(p.Hand))
			args.Value = p.Hand[randIdx].Value
		} else {
			args.Value = "-1"
		}
		p.callForCards(args)
	} else {
		time.Sleep(300 * time.Millisecond)
	}
}

/**
 * Calling for Card Request
*/
func (p *Player) callForCards(args CardRequest) () {
	reply := CardRequestReply{}

	if !call("GoFishServer.CallForCards", args, &reply) {
	    fmt.Println("Fail to request for Cards")
	}
	fmt.Println(reply.Cards)
	// append any returned cards to your hand
	if len(reply.Cards) > 1 || reply.Cards[0].Value != "-1" {
		p.Hand = append(p.Hand, reply.Cards...)
	}
}

/**
 * Calling for Game Status
*/
func (p *Player) callForGameStatus () (GameStatusReply, error) {
    args := GameStatusargs{}
  	reply := GameStatusReply{}
    if !call("GoFishServer.GetStatusOfGame", args, &reply) {
        fmt.Println("Fail to receive the Game status")
	}
	return reply, nil
}

/**
 * Enter the Game
*/
func (p *Player) enterGame () {

	fmt.Println("Join the Game ...")
    args := CardRequest{}
	reply := CardRequestReply{}
	// Ask for Joining the Game
    if !call("GoFishServer.EnterGame", args, &reply) {
        fmt.Println("Fail to Enter the Game")
        return
    }
    p.ID = reply.ID
    fmt.Println("PLAYER-ID : [", p.ID, "] Join the Game ...")
}


func main () {
	p := Player{} /* Create a new player */
	p.enterGame() /* Player enters the game */

	if p.ID == -1 {
		fmt.Println("Exceeds the maximum # ...")
		return
	}

	var gameOver = false
	for !gameOver {
	    var Gs, _ = p.callForGameStatus() // Checking the Game Status

	    /* renew the players hands */
		p.Players = Gs.Players
		p.Hand = Gs.Players[p.ID].Hand
		p.Pair = Gs.Players[p.ID].Pair

		time.Sleep(2 * time.Second)

		if Gs.CurrentPlayerId == p.ID && Gs.Ready {
			fmt.Println( "PLAYER [", p.ID, "] : Calling for Cards ...")
			p.callForStart()
			pairList, _ := p.HandsUpdates()
			p.callForEnd(pairList)
		}

		gameOver = Gs.Finished
	}
}