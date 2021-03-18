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