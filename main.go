package main

import (
	// hclogadapter "github.com/sarthakvk/hex-app/adapters/hclog"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	// "time"

	// "github.com/hashicorp/raft"

	"github.com/hashicorp/raft"
	hclogadapter "github.com/sarthakvk/hex-app/adapters/hclog"
	raftadapter "github.com/sarthakvk/hex-app/adapters/raft_adapter"
	// redisadapters "github.com/sarthakvk/hex-app/adapters/redis"
	// filehandling "github.com/sarthakvk/hex-app/core/file_handling"
)

var (
	myAddr = flag.String("address", "localhost:50051", "TCP host+port for this node")
	raftId = flag.String("raft_id", "", "Node id used by Raft")

	raftDir       = flag.String("raft_data_dir", "data/", "Raft data dir")
	raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
	logger        = hclogadapter.GetLogger()
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

	rft, net_transport := raftadapter.NewRaft(*raftId, host, portInt, *raftBootstrap)
	rft.LeadershipTransfer()
	rft.Apply()
	if <-rft.LeaderCh() {
		out := rft.AddVoter(raft.ServerID("B"), raft.ServerAddress("localhost:8001"), 0, time.Duration(time.Duration.Microseconds(100)))
		_ = rft.AddVoter(raft.ServerID("C"), raft.ServerAddress("localhost:8002"), 0, time.Duration(time.Duration.Microseconds(100)))
		if out.Error() != nil {
			panic(out.Error())

		}

		for {
			<-net_transport.Consumer()
		}
	}

}
