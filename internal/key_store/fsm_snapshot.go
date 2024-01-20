package keystore

import (
	"encoding/json"

	raft_lib "github.com/hashicorp/raft"
)

// KeyStoreSnapshot holds snapshot data returned by FSM, which will be saved to the SnapshotStore
// which then later can be recoved.
type KeyStoreSnapshot struct {
	data []byte
}

func newKeyStoreSnapshot(data_map map[string]string) (*KeyStoreSnapshot, error) {
	data, err := json.Marshal(data_map)

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	snapStore := &KeyStoreSnapshot{data: data}

	return snapStore, nil
}

// Persist is an implementation to write Snapshot data so that it can be saved
func (snap *KeyStoreSnapshot) Persist(sink raft_lib.SnapshotSink) error {
	_, err := sink.Write(snap.data)
	if err != nil {
		logger.Error(err.Error())
		sink.Cancel()
		return err
	}

	sink.Close()

	return nil
}

// This implements Release, it will be called when we will be done with the snapshot
func (snap *KeyStoreSnapshot) Release() {}
