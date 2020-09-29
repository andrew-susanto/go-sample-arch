package user

import "github.com/joez-tkpd/go-sample-arch/entity"

type Resource interface {
	GetUserByIDDB(id int64) entity.User

	GetUserByIDCache(id int64) entity.User
	SetUserCache(user entity.User) error
}

type resource struct {
	DB    DBRepository
	Cache CacheRepository
}

func NewResource(db DBRepository, cache CacheRepository) resource {
	return resource{
		DB:    db,
		Cache: cache,
	}
}
