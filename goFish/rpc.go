package goFish

/* RPC Definition */
import "os"
import "strconv"

/* RPC Definition */
type Args struct {X int}

type Reply struct {
	PlayerNum int
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}

