package goFish

import (
	"fmt"
	"log"
)

import "net/rpc"

func goFishP() {
	args := Args{}
	for {
		/* declare a reply structure */
		reply := Reply{}
		/* send the RPC request wait for the reply */
		if !call("goFishS.DistributeCard", args, &reply) { break }

		if reply.PlayerNum != -1 {

		}
	}
}


//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil { log.Fatal("dialing:", err) }
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}