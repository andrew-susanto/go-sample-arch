package cacheservice

import (
	// golang package
	"context"
	"fmt"
	"time"

	// external package
	"github.com/redis/go-redis/v9"
)

// Cache is interface of cache client
type Cache interface {
	// Set Redis `SET key value [expiration]` command. Use expiration for `SETEx`-like behavior.
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) *redis.StatusCmd

	// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
	Get(ctx context.Context, key string) *redis.StringCmd
}

// InitCache is function to initialize cache connection
func InitCache(host string, port string, password string) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: password,
		DB:       0, // use default DB
	})

	return rdb
}
