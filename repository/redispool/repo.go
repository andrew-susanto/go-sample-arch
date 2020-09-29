package redispool

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Repository struct {
	Redis *redis.Pool
}

func NewRedisPool(addr string) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     1,
		MaxActive:   1,
		IdleTimeout: time.Duration(10) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
		Wait: true,
	}
	return pool
}

func NewResource(redisPool *redis.Pool) Repository {
	return Repository{
		Redis: redisPool,
	}
}
