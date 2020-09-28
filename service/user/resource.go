package user

import (
	"github.com/joez-tkpd/go-sample-arch/entity"
	"github.com/joez-tkpd/go-sample-arch/repository/pgsqlx"
	"github.com/joez-tkpd/go-sample-arch/repository/redispool"
)

//go:generate mockgen -source=./resource.go -destination=./resource_mock.go -package=user

type Resource interface {
	GetUserByIDPgSqlx(id int64) entity.User

	GetUserByIDRedis(id int64) entity.User
	SetUserRedis(user entity.User) error
}

type resource struct {
	PgSqlx    pgsqlx.Repository
	RedisPool redispool.Repository
}

func NewResource(pgSqlx pgsqlx.Repository, redis redispool.Repository) Resource {
	return resource{
		PgSqlx:    pgSqlx,
		RedisPool: redis,
	}
}

func (rsc resource) GetUserByIDRedis(id int64) entity.User {
	user := rsc.RedisPool.GetUserByID(id)

	// convert to general entity object
	result := entity.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    user.Gender,
	}

	return result
}

func (rsc resource) GetUserByIDPgSqlx(id int64) entity.User {
	user := rsc.PgSqlx.GetUserByID(id)

	// convert to general entity object
	result := entity.User{
		ID:        user.ID,
		FirstName: user.Name,
		LastName:  "", // not provided
		Gender:    user.Gender,
	}

	return result
}

func (rsc resource) SetUserRedis(user entity.User) error {
	rsc.RedisPool.SetUser(redispool.User{
		user.ID,
		user.FirstName,
		user.LastName,
		user.Gender,
	})

	return nil
}
