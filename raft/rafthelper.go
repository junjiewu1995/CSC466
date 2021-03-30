package raft

import "time"
import "math/rand"


/**
 * Election time sets up
 * The Election time should be large than the HeartBeat Time according to the paper
*/

const HeartBeat = 50 * time.Millisecond
const MinElection = HeartBeat * 10
const MaxElection = MinElection * 8 / 5

/**
 * electionTimeout sets up
*/

func electionTimeout() time.Duration {
    return time.Duration(int(MinElection) + rand.Intn(int(MaxElection - MinElection)))
}

