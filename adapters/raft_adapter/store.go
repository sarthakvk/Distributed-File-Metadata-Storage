package raft

import (
	"fmt"
	"math/rand"

	raft_lib "github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

var (
	LogStoreFilePath     = "logstore.db" + fmt.Sprintf("%d", rand.Int())
	StableStoreFilePath  = "stablestore.db" + fmt.Sprintf("%d", rand.Int())
	SnapshotStoreBaseDir = "snapshot_data/" + fmt.Sprintf("%d", rand.Int())

	LogStore         = NewLogStore()
	StableStore      = NewStableStore()
	SnapshotStore, _ = raft_lib.NewFileSnapshotStore(SnapshotStoreBaseDir, 3, nil)
)

func NewStableStore() *raftboltdb.BoltStore {
	store, err := raftboltdb.NewBoltStore(StableStoreFilePath)

	if err == nil {
		return store
	}
	logger.Error("Failed to initialize StableStore:\n" + err.Error())

	return nil
}

func NewLogStore() *raftboltdb.BoltStore {
	store, err := raftboltdb.NewBoltStore(LogStoreFilePath)

	if err == nil {
		return store
	}

	logger.Error("Failed to initialize LogStore", err)

	return nil
}
