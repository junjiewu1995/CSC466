package raft

import (
	"fmt"
	"time"
	"../labrpc"
)

/**
 * // create a new Raft server instance:
 * rf := Make(peers, me, persister, applyCh)
 *
 * // start agreement on a new log entry:
 * rf.Start(command interface{}) (index, term, isLeader)
 *
 * // ask a Raft for its current term, and whether it thinks it is leader
 * rf.GetState() (term, isLeader)
 *
 * // each time a new entry is committed to the log, each Raft peer
 * // should send an ApplyMsg to the service (or tester).
 * type ApplyMsg
 *
 */

// Make is
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {

	rf := &Raft{}

	// Make 函数参数的去处
	rf.peers = peers
	rf.me = me
	rf.persister = persister
	rf.chanApply = applyCh

	// 需要 persist 的参数
	rf.currentTerm = 0
	rf.votedFor = NOBODY
	le := LogEntry{LogIndex: 0, LogTerm: 0, Command: 0}
	rf.logs = append(rf.logs, le) // 在 logs 预先放入一个，方便 Raft.getLastIndex()

	// 私有状态
	rf.state = FOLLOWER
	
    // REVIEW: 把这些通道的设置成非缓冲的，看看会不会出错
	rf.chanCommit = make(chan struct{}, 100)
	rf.chanHeartBeat = make(chan struct{}, 100)
	rf.chanBeElected = make(chan struct{}, 100)

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	go rf.statesLoop()

	go rf.applyLoop()

	return rf
}

func (rf *Raft) statesLoop() {
	for {
		switch rf.state {
		case FOLLOWER:
			select {
			case <-time.After(electionTimeout()): 
                rf.state = CANDIDATE
			case <-rf.chanHeartBeat:
			}
		case CANDIDATE:
			rf.newElection()
		case LEADER:
			rf.newHeartBeat()
		}
	}
}

func (rf *Raft) newElection() {
	rf.mu.Lock()

	rf.currentTerm++
	rf.votedFor = rf.me
	rf.voteCount = 1

	rf.persist()
	rf.mu.Unlock()

	go rf.broadcastRequestVote()

	select {
		case <-time.After(electionTimeout()):
		case <-rf.chanHeartBeat:
			rf.state = FOLLOWER
		case <-rf.chanBeElected:
			rf.comeToPower()
	}
}

func (rf *Raft) comeToPower() {
	
    rf.mu.Lock()
	rf.state = LEADER
	fmt.Printf("%s is Leader now\n", rf)
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))
	for i := range rf.peers {
		rf.nextIndex[i] = rf.getLastIndex() + 1
		rf.matchIndex[i] = 0
	}
	rf.mu.Unlock()
}

func (rf *Raft) newHeartBeat() {
	rf.broadcastAppendEntries()
	<-time.After(heartBeat)
}

func (rf *Raft) applyLoop() {
	for {
		<-rf.chanCommit
		rf.mu.Lock()

		commitIndex := rf.commitIndex
		baseIndex := rf.getBaseIndex()
        
		for i := rf.lastApplied + 1; i <= commitIndex; i++ {
            
			msg := ApplyMsg {
				CommandValid: true,
				CommandIndex: i,
				Command: rf.logs[i-baseIndex].Command,
			}
			rf.chanApply <- msg
			rf.lastApplied = i
		}
		rf.mu.Unlock()
	}
}

// Start is
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and
// ** return immediately, without waiting for the log appends to complete. **
// there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// The first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. The third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
    
	// Your code here (2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if !rf.isLeader() {	return -1, -1, false }
	logIndex := rf.getLastIndex() + 1
	term := rf.currentTerm
	isLeader := rf.isLeader()

	rf.logs = append(rf.logs,
		LogEntry{
			LogIndex: logIndex,
			LogTerm:  term,
			Command:  command,
		}) // append new entry from client
	rf.persist()

	// Your code above
	return logIndex, term, isLeader
}

// GetState is
// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isLeader bool
	// Your code here (2A).
	term = rf.currentTerm
	isLeader = rf.isLeader()

	// Your code above (2A)
	return term, isLeader
}

// ApplyMsg is
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool // CommandValid = true 表示， 此消息是用于应用 
	CommandIndex int  // Command 所在的 logEntry.logIndex 值
	Command      interface{}
    
}

func (m ApplyMsg) String() string {
	return fmt.Sprintf("ApplyMsg{Valid:%t,Index:%d,Command:%v}", m.CommandValid, m.CommandIndex, m.Command)
}

// Kill is
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}
