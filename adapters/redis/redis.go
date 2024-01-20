package redisadapters

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sarthakvk/gofilemeta/adapters/logging"
)

var logger = logging.GetLogger()

// RedisBacked holds the connection to the redis server
type RedisBackend struct {
	rdb *redis.Client
}

// Creates a new redis backend
func NewRedisBackend(address string) *RedisBackend {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rfsm := RedisBackend{rdb: rdb}

	return &rfsm
}

// Get retrieves the value associated with the given key from Redis.
func (rb *RedisBackend) Get(key string) (string, bool) {
	value, err := rb.rdb.Get(context.Background(), key).Result()
	if err != nil {
		// Key does not exist
		return "", false
	}

	return value, true
}

// Delete deletes the key and its associated value from Redis.
func (rb *RedisBackend) Delete(key string) error {
	err := rb.rdb.Del(context.Background(), key).Err()
	if err == redis.Nil {
		// Key does not exist
		return fmt.Errorf("key '%s' not found", key)
	} else if err != nil {
		// Other errors
		return err
	}

	return nil
}

// Set sets the value with the key inside redis
func (r *RedisBackend) Set(key, value string) error {
	var ctx = context.Background()
	_, err := r.rdb.Set(ctx, key, value, time.Duration(0)).Result()

	if err != nil {
		logger.Error(err.Error())
	}

	return err
}

// Snapshot, it implements Snapshot to support backup
func (r *RedisBackend) Snapshot() map[string]string {
	data := make(map[string]string)
	keys, err := r.rdb.Keys(context.Background(), "*").Result()
	if err != nil {
		logger.Error(err.Error())
	}

	for _, key := range keys {
		value, err := r.rdb.Get(context.Background(), key).Result()
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		data[key] = value

	}
	return data
}

// It implements Restore, so the node can recover from failures
func (r *RedisBackend) Restore(data map[string]string) error {

	logger.Debug("[START RESTORE] read all message from snapshot")
	var totalRestored int

	for key, value := range data {

		r.Set(key, value)
		totalRestored++
	}

	logger.Debug("[END RESTORE] success restore %d messages in snapshot\n", totalRestored)
	return nil
}
