package redispool

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=redispool

type Repository interface {
	SetUser(user User)
	GetUserByID(id int64) User
}

type repository struct {
	Redis *redis.Pool
}

func NewRedisPool() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     1,
		MaxActive:   1,
		IdleTimeout: time.Duration(10) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
		Wait: true,
	}
	return pool
}

func NewResource(redisPool *redis.Pool) Repository {
	return repository{
		Redis: redisPool,
	}
}

func (repo repository) GetUserByID(id int64) User {
	strID := strconv.FormatInt(id, 10)

	conn := repo.Redis.Get()
	defer conn.Close()

	resp, _ := redis.String(conn.Do("GET", "user_"+strID))

	user := new(User)
	json.Unmarshal([]byte(resp), &user)

	return *user
}

func (repo repository) SetUser(user User) {
	strID := strconv.FormatInt(user.ID, 10)

	conn := repo.Redis.Get()
	defer conn.Close()

	encoded, _ := json.Marshal(user)
	conn.Do("SET", "user_"+strID, encoded)

	return
}
