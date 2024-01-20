package keystore

import (
	"encoding/json"

	raft_lib "github.com/hashicorp/raft"
)

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

func (snap *KeyStoreSnapshot) Release() {}
