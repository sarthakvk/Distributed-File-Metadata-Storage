package raft

import (
	"time"

	raft_lib "github.com/hashicorp/raft"

	"github.com/sarthakvk/gofilemeta/adapters/logging"
)

const (
	ElectionWaitDuration    = time.Duration(time.Millisecond * 10)
	AddVoterTimeoutDuration = time.Duration(time.Millisecond * 500)
	ApplyTimeout            = time.Duration(time.Second)
)

var logger = logging.GetLogger()

// Raft defines an adapter which will be used
// to communicate with the raft cluster
type Raft struct {
	raftID    raft_lib.ServerID
	node      *raft_lib.Raft
	transport *raft_lib.NetworkTransport
}

func (r *Raft) IsLeader() bool {
	return r.node.State() == raft_lib.Leader
}

// GetLeader will return the Address of the current leader
// If there is no leader, it will wait till new leader is elected.
func (r *Raft) GetLeader() raft_lib.ServerAddress {
	leader := raft_lib.ServerAddress("")
	for {
		leader, _ = r.node.LeaderWithID()
		if leader == "" {
			time.Sleep(ElectionWaitDuration)
		} else {
			break
		}
	}
	return leader
}

// Creates a new raft node
// TODO : Currently using default config, provide a way to configure the Raft node
func NewRaft(raftID, address string, fsm raft_lib.FSM) *Raft {
	config := raft_lib.DefaultConfig()
	config.LocalID = raft_lib.ServerID(raftID)

	transport := NewTransport(address)

	node, err := raft_lib.NewRaft(config, fsm, LogStore, StableStore, SnapshotStore, transport)

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	return &Raft{raftID: raft_lib.ServerID(raftID), node: node, transport: transport}
}

// Bootstrap Cluster
func (r *Raft) BootstrapCluster() {
	config := raft_lib.Configuration{
		Servers: []raft_lib.Server{{
			Suffrage: raft_lib.Voter,
			Address:  r.transport.LocalAddr(),
			ID:       r.raftID,
		}},
	}

	fut := r.node.BootstrapCluster(config)
	err := fut.Error()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}

// Add voter to our cluster, Note:L this must be called from a leader
// if the caller is not a leader, then it fails silently logging the error
func (r *Raft) AddVoter(raftID, address string) {
	serverID := raft_lib.ServerID(raftID)
	serverAddress := raft_lib.ServerAddress(address)
	prevIndex := uint64(0)

	future := r.node.AddVoter(serverID, serverAddress, prevIndex, AddVoterTimeoutDuration)
	err := future.Error()

	if err != nil {
		logger.Debug("Failed to add voter raft node!")
		logger.Error(err.Error())
	}

}

// Apply is used to apply a command to the FSM in a highly consistent
// manner. This returns a future that can be used to wait on the application.
// This must be run on the leader or it will fail.
func (r *Raft) Apply(cmd []byte) raft_lib.ApplyFuture {
	return r.node.Apply(cmd, ApplyTimeout)
}
