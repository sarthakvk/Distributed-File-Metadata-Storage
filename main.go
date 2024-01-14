package main

import (
	"flag"
	"fmt"
	raftadapter "github.com/sarthakvk/hex-app/adapters/raft_adapter"
	"net"
	"strconv"
)

var (
	myAddr        = flag.String("address", "localhost:50051", "TCP host+port for this node")
	raftId        = flag.String("raft_id", "", "Node id used by Raft")
	raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
)

func main() {
	flag.Parse()

	if *raftId == "" {
		panic("flag --raft_id is required")
	}

	host, port, err := net.SplitHostPort(*myAddr)
	portInt, _ := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("failed to parse local address (%q): %v", *myAddr, err))
	}

	_, net_transport := raftadapter.NewRaft(*raftId, host, portInt, *raftBootstrap)

	for {
		<-net_transport.Consumer()
	}

}
