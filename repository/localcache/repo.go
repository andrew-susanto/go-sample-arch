package localcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=localcache

func NewRepository(gocache *cache.Cache) Repository {
	return repository{
		GoCache: gocache,
	}
}

type repository struct {
	GoCache *cache.Cache
}

func NewGoCache() *cache.Cache {
	return cache.New(15*time.Minute, 10*time.Minute)
}

// type Repository interface {
// 	Get(string) (val interface{}, isExists bool)
// 	GetString(string) (val string, isExists bool)
// 	SetNoExpiry(string, interface{}) error
// 	Delete(string) error
// }

func (repo repository) Get(key string) (interface{}, bool) {
	return repo.GoCache.Get(key)
}

func (repo repository) GetString(key string) (string, bool) {
	val, isExists := repo.GoCache.Get(key)
	strVal, _ := val.(string)
	return strVal, isExists
}

func (repo repository) SetNoExpiry(key string, val interface{}) error {
	repo.GoCache.Set(key, val, cache.NoExpiration)
	return nil
}

func (repo repository) Delete(key string) error {
	repo.GoCache.Delete(key)
	return nil
}
