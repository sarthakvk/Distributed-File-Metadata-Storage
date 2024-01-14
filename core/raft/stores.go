package raft_integration

import (
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	hclogadapter "github.com/sarthakvk/hex-app/adapters/hclog"
)

const LogStoreFilePath string = "logstore.db"
const StableStoreFilePath string = "stablestore.db"
const SnapshotStoreBaseDir = "snapshot_data/"

var (
	LogStore         = NewLogStore()
	StableStore      = NewStableStore()
	SnapshotStore, _ = raft.NewFileSnapshotStore(SnapshotStoreBaseDir, 3, nil)
)

func NewStableStore() *raftboltdb.BoltStore {
	store, err := raftboltdb.NewBoltStore(StableStoreFilePath)

	if err == nil {
		return store
	}

	hclogadapter.Logger.Exception("Failed to initialize StableStore", err)

	return nil
}

func NewLogStore() *raftboltdb.BoltStore {
	store, err := raftboltdb.NewBoltStore(LogStoreFilePath)

	if err == nil {
		return store
	}

	hclogadapter.Logger.Exception("Failed to initialize LogStore", err)

	return nil
}
