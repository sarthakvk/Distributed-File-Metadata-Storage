package raft

import (
	"fmt"
	"net"
	"time"

	raft_lib "github.com/hashicorp/raft"
)

func SplitHostPort(address string) (string, string) {
	host, port, err := net.SplitHostPort(address)

	if err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}
	return host, port
}

func NewTransport(address string) *raft_lib.NetworkTransport {
	maxPool := 5
	timeout := time.Duration(time.Duration.Milliseconds(100))

	host, _ := SplitHostPort(address)

	ips, err := net.LookupIP(host)

	if err != nil || len(ips) == 0 {
		panic(fmt.Errorf("can't resolve host: %s\n%s", host, err.Error()))
	}

	advertise, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		panic(err)
	}
	tran, err := raft_lib.NewTCPTransportWithLogger(address, advertise, maxPool, timeout, logger)

	if err != nil {
		panic(err)
	}

	return tran
}
