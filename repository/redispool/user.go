package redispool

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

const (
	keyUser = "user_"
)

func (repo Repository) GetUserByID(id int64) User {
	strID := strconv.FormatInt(id, 10)

	conn := repo.Redis.Get()
	defer conn.Close()

	resp, err := redis.String(conn.Do("GET", keyUser+strID))
	log.Println(err)

	user := new(User)
	json.Unmarshal([]byte(resp), &user)

	return *user
}

func (repo Repository) SetUser(user User) {
	strID := strconv.FormatInt(user.ID, 10)

	conn := repo.Redis.Get()
	defer conn.Close()

	encoded, _ := json.Marshal(user)
	conn.Do("SET", keyUser+strID, encoded)

	return
}
