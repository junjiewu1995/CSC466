package main

import (
	"fmt"
	"log"
	"net/rpc"
)

import . "../CallingUtilities"

/**
 * Cient maintains its own Hands
*/

type Player struct {
	ID 			int      // Player ID
	Hand 		[]Card   // Current Card in Hand
	Pairs 		[]Pairs  // Pairs for wins
	rivals	    []Player // Other Players
	Win         bool     // Check for Game Over
}

type GoFishGameReply struct {
	Turn 		int
}

func (p *Player ) gameOver () bool { return p.Win }

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


/**
 * Players Hands Updates
*/
func (p *Player) HandsUpdates (C Card, Cards []Card) error {



    return nil
}

/**
 * Calling for Card Request
*/

func (p *Player) CallCardRequest(goFish bool) {
	args := CardRequest{}
	reply := CardRequestReply{}

    if goFish {
        args.goFish = true
    }

    // Ask for a random Card from a random Players

    p.HandsUpdates()

	// Ask for a Card Components
	if !call("GoFishServer.RequestForCard", args, &reply) {
	    fmt.Println("Fail to request for Cards")
	    return
	}
	return reply
}


/**
 * Calling for Game Status
*/

func (p *Player) GameStatus () GameStatusReply {

    args := GameStatusRequest{}
    gameStatusreply := GameStatusReply{}

    // Ask for a Card Components
    if !call("GoFishServer.GetStatusOfGame", args, &gameStatusreply) {
        fmt.Println("Fail to receive the Game status")
        return
    }
    return gameStatusreply
}

/**
 * Enter the Game
*/

func (p *Player) EnterGame () {
    args := CardRequest{}
	reply := CardRequestReply{}

	// Ask for Joining the Game
    if !call("GoFishServer.EnterGame", args, &reply) {
        fmt.Println("Fail to Enter the Game")
        return
    }
    fmt.Println(reply.ID)
}


func main () {
	gsc := Player{}
	gsc.EnterGame()

	for !gsc.gameOver() {
	    gameStatusreply := GameStatusReply{}

	    /* Checking the Game Status */
	    gameStatusreply = gsc.GameStatus()

	    if gameStatusreply.Game {
	        p.Win = true
	    }

	    if gameStatusreply.Turn == gsc.ID {
	        gsc.CardRequest(gameStatusreply.goFish)
	    }

	    /* Check the Game Status for 2 secs */
        time.Sleep( 2 * time.Second)
	}
}