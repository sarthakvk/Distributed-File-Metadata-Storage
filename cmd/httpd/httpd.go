package main

import (
	"flag"

	redis "github.com/sarthakvk/gofilemeta/adapters/redis"
	"github.com/sarthakvk/gofilemeta/internal/httpd"
	keystore "github.com/sarthakvk/gofilemeta/internal/key_store"
)

var (
	bootstrap     = flag.Bool("bootstrap", false, "Bootstrap Cluster")
	addr          = flag.String("address", "localhost:8000", "address of the raft node")
	nodeID        = flag.String("node-id", "", "unique id for node")
	httpd_port    = flag.Int("http-port", 9000, "Port to run running HTTP service")
	redis_address = flag.String("redis-addr", "localhost:6379", "address of the redis server")
)

func main() {
	flag.Parse()
	dataBackend := redis.NewRedisBackend(*redis_address)
	store := keystore.New(*nodeID, *addr, dataBackend, *bootstrap)

	httpd.RunServer(store, *httpd_port)
}
