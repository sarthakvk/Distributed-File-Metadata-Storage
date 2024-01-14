package raft_integration

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
)

func NewTCPTransport(host string, port int) *raft.NetworkTransport {
	maxPool := 5
	logOutput := hclog.New(&hclog.LoggerOptions{Name: "network-logs"})
	timeout := time.Duration(time.Duration.Milliseconds(100))

	ips, err := net.LookupIP(host)

	if err != nil || len(ips) == 0 {
		panic(fmt.Errorf("Can't resolve host: %s\n%s", host, err.Error()))
	}

	advertise := &net.TCPAddr{IP: ips[0], Port: port}

	tran, _ := raft.NewTCPTransportWithLogger("0.0.0.0:0", advertise, maxPool, timeout, logOutput)

	return tran
}
