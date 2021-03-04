package goFish

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type goFishS struct {
	mux sync.Mutex
	NPlayer int
}

func (m *goFishS) DistributeCard() {
	fmt.Println("Hi Player")
}

func (m *goFishS) server() {
	//start a thread that listens for RPCs from goFish Player.go
	rpc.Register(m)
	rpc.HandleHTTP()
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {log.Fatal("listen error:", e)}
	go http.Serve(l, nil)
}