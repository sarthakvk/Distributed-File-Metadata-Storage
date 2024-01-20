package keystore_adapter

// AbstractKeyStore define the interface which can
// be used as a key-value pair storage
type AbstractKeyStore interface {
	Get(key string) (string, bool)
	Set(key, value string) error
	Delete(key string) error
	GetOrCreate(key, value string) (bool, string, error)
	Replicate(nodeID, address string)
}
