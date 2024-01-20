package keystore_adapter

type AbstractKeyStore interface {
	Get(key string) (string, bool)
	Set(key, value string) error
	Delete(key string) error
	GetOrCreate(key, value string) (bool, string, error)
	Replicate(nodeID, address string)
}
