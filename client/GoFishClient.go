package main

import (
	"fmt"
	"log"
	"net/rpc"
)
import . "../CallingUtilities"

type Player struct {
	ID 			int
	Hand 		[]Card
	Pairs 		[]Pairs
	Opponents 	[]Player
}

type GoFishGameReply struct {
	Turn 		int
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

func (p *Player)goFishP() {
	fmt.Println("GoFishPlayer")
}

func (p *Player) CallCardRequest() {
	args := CardRequest{}
	reply := CardRequestReply{}

	// Ask for a Card Components
	if !call("GoFishServer.RequestForCard", args, &reply) {
	    fmt.Println("Hello World :)")
	    return
	}

}

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
	gsc.goFishP()
	gsc.EnterGame()
	gsc.CallCardRequest()
}

