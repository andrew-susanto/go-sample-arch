package cache

import (
	// golang package
	"context"
	"time"

	// external package
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=cache

type Redis interface {
	// Set Redis `SET key value [expiration]` command. Use expiration for `SETEx`-like behavior.
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) *redis.StatusCmd

	// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
	Get(ctx context.Context, key string) *redis.StringCmd
}

type Repository struct {
	redis Redis
}

// NewRepository is a function to initialize cache repository
func NewRepository(redis Redis) Repository {
	return Repository{
		redis: redis,
	}
}
