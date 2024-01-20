package main

import (
	"flag"

	redis "github.com/sarthakvk/hex-app/adapters/redis"
	"github.com/sarthakvk/hex-app/internal/httpd"
	keystore "github.com/sarthakvk/hex-app/internal/key_store"
)

var (
	bootstrap  = flag.Bool("bootstrap", false, "Bootstrap Cluster")
	addr       = flag.String("address", "localhost:8000", "address of the raft node")
	nodeID     = flag.String("node-id", "", "unique id for node")
	httpd_port = flag.Int("port", 9000, "Port to run running HTTP service")
)

func main() {
	flag.Parse()
	dataBackend := redis.NewRedisBackend()
	store := keystore.New(*nodeID, *addr, dataBackend, *bootstrap)

	httpd.RunServer(store, *httpd_port)
}
