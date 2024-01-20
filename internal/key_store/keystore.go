package keystore

import (
	"errors"
	"sync"

	ks_adapter "github.com/sarthakvk/gofilemeta/adapters/keystore_adapter"
	"github.com/sarthakvk/gofilemeta/adapters/logging"
	raft "github.com/sarthakvk/gofilemeta/adapters/raft_adapter"
)

var (
	logger = logging.GetLogger()
)

// DataBackend will be implemented by
// the backend used for storage
type DataBackend interface {
	Get(string) (string, bool)
	Set(string, string) error
	Delete(string) error
	Snapshot() map[string]string
	Restore(map[string]string) error
}

// KeyStore is the distributed Key-value store
// It uses `Raft` internally to get consensus from the peers
type KeyStore struct {
	// Consesus mechanism
	raft *raft.Raft

	// Data storage for the key-value store
	data DataBackend

	rw_lock sync.RWMutex
}

// Create new KeyStore.
// bootstrap_cluster is flag; whether to create a Raft cluster. This only needs to be ran once in the lifespan of a cluster
func New(nodeID, address string, backend DataBackend, bootstrap_cluster bool) *KeyStore {
	store := &KeyStore{}
	fsm := (*KeyStoreFSM)(store)
	store.data = backend
	raft := raft.NewRaft(nodeID, address, fsm)
	store.raft = raft

	if bootstrap_cluster {
		go raft.BootstrapCluster()
	}

	return store
}

// Add voters to the cluster, must be ran from leader otherwise it will fail
// Our Consesnsus will make sure to bring the replica to the latest state
func (store *KeyStore) Replicate(nodeID, address string) {
	go store.raft.AddVoter(nodeID, address)
}

// Get the Value from data backend
func (store *KeyStore) Get(key string) (string, bool) {
	store.rw_lock.Lock()
	defer store.rw_lock.Unlock()

	value, ok := store.data.Get(key)

	return value, ok
}

// Delete, deletes the given key in highly consistent manner
//
// It must be called through leader otherwise it will return an error
func (store *KeyStore) Delete(key string) error {
	if !store.raft.IsLeader() {
		logger.Error("key deletion request denied, not a leader")
		return errors.New("request denied, Not a leader")
	}

	var cmd ks_adapter.Command

	cmd.Operation = ks_adapter.DELETE
	cmd.Key = key

	raw_cmd, err := getRawCommand(cmd)

	if err != nil {
		logger.Error("failed to marshal the delete request command")
		return errors.New("failed to marshal the delete request command")
	}

	fut := store.raft.Apply(raw_cmd)

	if err := fut.Error(); err != nil {
		logger.Error(err.Error())
		return err
	}

	resp := fut.Response()

	switch resp := resp.(type) {
	case error:
		return resp
	default:
		return nil
	}

}

// Set, sets the given key with value in highly consistent manner
//
// It must be called through leader otherwise it will return an error
func (store *KeyStore) Set(key, value string) error {
	if !store.raft.IsLeader() {
		logger.Error("key deletion request denied, Not a leader!")
		return errors.New("request denied, not a leader")
	}

	var cmd ks_adapter.Command

	cmd.Operation = ks_adapter.SET
	cmd.Key = key
	cmd.Value = value

	raw_cmd, err := getRawCommand(cmd)

	if err != nil {
		logger.Error("failed to marshal the delete request command")
		return errors.New("failed to marshal the delete request command")
	}

	fut := store.raft.Apply(raw_cmd)

	if err := fut.Error(); err != nil {
		logger.Error(err.Error())
		return err
	}

	resp := fut.Response()

	switch resp := resp.(type) {
	case error:
		return resp
	default:
		return nil

	}
}

// GetOrCreate, gets the given key or creates incase it doesn't exist in highly consistent manner
//
// It must be called through leader otherwise it will return an error
// returns: created, value, error
func (store *KeyStore) GetOrCreate(key, value string) (bool, string, error) {
	created := false
	val, ok := store.Get(key)

	if ok {
		return created, val, nil
	}
	err := store.Set(key, value)

	if err == nil {
		created = true
	}

	return created, value, err

}
