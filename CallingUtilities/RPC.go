package helper

import (
	"os"
	"strconv"
)


type Player struct {
	ID 			int
	Hand 		[]Card
	Pairs 		[]Pairs
	Opponents 	[]Player
}

type Args struct {
	X int
}

type Reply struct {
	PlayerNum int
}

type Card struct {
	/* Value Range: 2,3,4,5,6,7,8,9,10,J,Q,K,A */
	/* Suit Range: clubs,diamonds,hearts,spades */
	Value string
	Suit string
	Used bool
}

type Pairs struct {
	One Card
	Two Card
}

type CardRequestReply struct {
	Turn 			int
	GoFishGame		bool
	ID              int
}

type CardRequest struct {
	GoFishGame 		bool
	ID              int
	P               Player
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
