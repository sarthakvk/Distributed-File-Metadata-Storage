package keystore

import (
	"encoding/json"
	"io"

	raft_lib "github.com/hashicorp/raft"
	ks_adapter "github.com/sarthakvk/hex-app/adapters/keystore_adapter"
)

type KeyStoreFSM KeyStore

func (store *KeyStoreFSM) Apply(log *raft_lib.Log) interface{} {

	cmd, err := GetCommand(log.Data)

	if err != nil {
		return err
	}

	switch cmd.Operation {
	case ks_adapter.SET:
		store.applySet(cmd.Key, cmd.Value)
	case ks_adapter.DELETE:
		store.applyDel(cmd.Key)
	}

	return nil
}

func (store *KeyStoreFSM) Snapshot() (raft_lib.FSMSnapshot, error) {
	store.rw_lock.RLock()
	defer store.rw_lock.RUnlock()

	data := store.data.Snapshot()
	snap, err := newKeyStoreSnapshot(data)

	if err != nil {
		logger.Error("Failed to create snapshot")
		return nil, err
	}
	return snap, nil
}

func (store *KeyStoreFSM) Restore(snapshot io.ReadCloser) error {
	store.rw_lock.Lock()
	defer store.rw_lock.Unlock()

	var (
		snap_data     []byte
		restored_data map[string]string
	)

	snapshot.Read(snap_data)

	err := json.Unmarshal(snap_data, &restored_data)
	if err != nil {
		logger.Error("Encountered erros while restoring snapshot data")
		return err
	}

	store.data.Restore(restored_data)

	return nil
}

func (store *KeyStoreFSM) applySet(key, value string) {
	store.rw_lock.Lock()
	defer store.rw_lock.Unlock()

	store.data.Set(key, value)
}

func (store *KeyStoreFSM) applyDel(key string) {
	store.rw_lock.Lock()
	defer store.rw_lock.Unlock()

	store.data.Delete(key)
}
