package main

import (
	"flag"
	"fmt"

	keystore "github.com/sarthakvk/gokey/internal/key_store"
)

var (
	bootstrap = flag.Bool("bootstrap", false, "Bootstrap Cluster")
	addr      = flag.String("address", "", "address of the node")
	raftID    = flag.String("node-id", "", "unique id for node")
	peerCount = flag.Int("peer-count", 0, "Number of peers to add; this only makes sense when bootstraping a cluster")
)

func main() {
	flag.Parse()

	store := keystore.New(*raftID, *addr, *bootstrap)

	if *bootstrap {
		addPeers(store, *peerCount)
	}

}

func addPeers(store *keystore.KeyStore, count int) {
	for i := 0; i < count; i++ {
		var raftID, address string
		println("Enter Peer: ", i+1)
		fmt.Scan(&raftID, &address)
		store.AddVoterNodes(raftID, address)
	}
}
