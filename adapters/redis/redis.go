package redisadapters

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/redis/go-redis/v9"
	// raft_integration "github.com/sarthakvk/hex-app/core/raft"
)

type Payload struct {
	Key   string
	Value string
}

type RedisFSM struct {
	rdb *redis.Client
}

func NewRedisFSM() *RedisFSM {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rfsm := RedisFSM{rdb: rdb}

	return &rfsm
}

func (r RedisFSM) set(key, value string) {
	var ctx = context.Background()
	r.rdb.Set(ctx, key, value, time.Duration(0))
}

func (r RedisFSM) Apply(log *raft.Log) interface{} {
	if log.Type == raft.LogCommand {
		payload := Payload{}
		if err := json.Unmarshal(log.Data, &payload); err != nil {
			panic(err)
		}
		r.set(payload.Key, payload.Value)
	}
	return nil
}

func (r RedisFSM) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshotNoop()
}

// Restore is used to restore an FSM from a Snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
// Restore will update all data in BadgerDB
func (r RedisFSM) Restore(rClose io.ReadCloser) error {
	defer func() {
		if err := rClose.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[FINALLY RESTORE] close error %s\n", err.Error())
		}
	}()

	_, _ = fmt.Fprintf(os.Stdout, "[START RESTORE] read all message from snapshot\n")
	var totalRestored int

	decoder := json.NewDecoder(rClose)
	for decoder.More() {
		var data = &Payload{}
		err := decoder.Decode(data)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error decode data %s\n", err.Error())
			return err
		}

		r.set(data.Key, data.Value)
		totalRestored++
	}

	// read closing bracket
	_, err := decoder.Token()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error %s\n", err.Error())
		return err
	}

	_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] success restore %d messages in snapshot\n", totalRestored)
	return nil
}
