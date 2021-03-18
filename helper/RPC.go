package helper

import (
	"os"
	"strconv"
)

type Args struct {
	X int
}

type Card struct {
	Value string
	Suit string
	Used bool
}

type Pairs struct {
	One Card
	Two Card
}

type Reply struct {
	PlayerNum int
}

type CardRequestReply struct {
	Turn 			int
	GoFishGame		bool
}

type CardRequest struct {
	GoFishGame 		bool
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
