package main

import (
	. "../helper"
	"fmt"
	"log"
	"net/rpc"
)

func goFishP() {
	fmt.Println("ABC")
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

func main () {
	fmt.Println("ABC")
}

