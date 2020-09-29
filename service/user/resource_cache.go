package user

import (
	"github.com/joez-tkpd/go-sample-arch/entity"
	"github.com/joez-tkpd/go-sample-arch/repository/redispool"
)

//go:generate mockgen -source=./resource_cache.go -destination=./resource_cache_mock.go -package=user

type CacheRepository interface {
	SetUser(user redispool.User)
	GetUserByID(id int64) redispool.User
}

func (rsc resource) GetUserByIDCache(id int64) entity.User {
	user := rsc.Cache.GetUserByID(id)

	// convert to general entity object
	result := entity.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    user.Gender,
	}

	return result
}

func (rsc resource) SetUserCache(user entity.User) error {
	rsc.Cache.SetUser(redispool.User{
		user.ID,
		user.FirstName,
		user.LastName,
		user.Gender,
	})

	return nil
}
