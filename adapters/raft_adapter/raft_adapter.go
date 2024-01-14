package raftadapter

import (
	"fmt"
	"github.com/hashicorp/raft"
	hclogadapter "github.com/sarthakvk/hex-app/adapters/hclog"
	redisadapters "github.com/sarthakvk/hex-app/adapters/redis"
	raft_integration "github.com/sarthakvk/hex-app/core/raft"
)

type RaftAdapter struct {
}

func NewRaftAdapter() RaftAdapter {
	return RaftAdapter{}
}

func (raft *RaftAdapter) GetConsensus() bool {
	hclogadapter.Logger.Log("Getting Consesnsus through RAFT")
	return true
}

func NewRaft(id, host string, port int, bootstrap bool) (*raft.Raft, *raft.NetworkTransport) {
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(id)
	rFsm := redisadapters.NewRedisFSM()
	trans := raft_integration.NewTCPTransport(host, port)
	node, err := raft.NewRaft(conf, rFsm, raft_integration.LogStore, raft_integration.StableStore, raft_integration.SnapshotStore, trans)

	if err != nil {
		hclogadapter.Logger.Exception("Failed to create raft node", err)
		panic(err)
	}

	if bootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{{
				Suffrage: raft.Voter,
				Address:  raft.ServerAddress(fmt.Sprintf("%s:%d", host, port)),
				ID:       raft.ServerID(id),
			}},
		}
		node.BootstrapCluster(cfg)
	}
	return node, trans
}
