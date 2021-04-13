package helper

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)


type Args struct {
	X int
}

type GameStateArgs struct {
	Key     string
	Payload string
	Ok      bool
}
type GameStateReply struct {
	Payload string
	Ok      bool
}

type Player struct {
	ID 			int
	Hand 		[]Card
	Pair 		[]Pairs
	Players 	[]Player
}


type Reply struct {
	PlayerNum int
}

func RaftServerSock() string {
	s := "/var/tmp/824-rb-"
	s += strconv.Itoa(os.Getuid())
	return s
}

type PlayPairRequest struct {
	Turn  int
	Owner int
	Hand  []Card
	Pair  []Pairs
}

type PlayPairReply struct {
	Turn     int
	Accepted bool
}


/* Value Range: 2,3,4,5,6,7,8,9,10,J,Q,K,A */
/* Suit Range: clubs,diamonds,hearts,spades */
type Card struct {
	Value string
	Suit string
	Used bool
}

type Pairs struct {
	One 	Card
	Two 	Card
	Three   Card
	Four    Card
}

type CardRequestReply struct {
	Turn 			int
	goFish  		bool
	ID              int
    Cards           []Card
	GoFish 			bool
}

type CardRequest struct {
	goFish   		bool
	PlayerID 		int
	GameTurn 		int
	PInfo   		Player
	MCard    		string

	Turn   			int
	Target 			int //Index of target player
	Value  			string
}



// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets
func MasterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}

// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func CallRB(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := RaftServerSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}