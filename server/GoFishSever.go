package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

import . "../helper"


type Player struct {
	ID 			int
	Hand 		[]Card
	Pairs 		[]Pairs
	Opponents 	[]Player
}

type GoFishServer struct {
	Mu 				sync.Mutex
	PlayerCounter 	int
}

func (m *GoFishServer) DistributeCard() {
	fmt.Println("~ Server Starts to Run ~")
}

/**
 * RPC server interaction
*/
func (m *GoFishServer) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	sockname := MasterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

/* Create a Game Server */
func StartGoFishServer () *GoFishServer {
	g := GoFishServer{}
	g.DistributeCard()
	return nil
}

/* Main Function */
func main () {
	/* Start to Run Go Fish Server */
	StartGoFishServer()
	/* Go Fish Server */
	fmt.Println("Game Over")
}