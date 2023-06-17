package user

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/repository/cache"
	"github.com/andrew-susanto/go-sample-arch/repository/psql"
)

//go:generate mockgen -source=./resource.go -destination=./resource_mock.go -package=user

type CacheRepository interface {
	// GetUserByIDFromCache gets user by id from cache service
	//
	// Returns entity user and nil error if success
	// Otherwise returns empty entity user and non-nil error
	GetUserByID(ctx context.Context, ID int64) (cache.User, error)
}

type DBRepository interface {
	// GetUserByIDFromDB gets user by id from db service
	//
	// Returns entity user and nil error if success
	// Otherwise returns empty entity user and non-nil error
	GetUserByID(ctx context.Context, ID int64) (psql.User, error)
}

type resource struct {
	cache CacheRepository
	db    DBRepository
}

// NewResource creates new user resource
func NewResource(cache CacheRepository, db DBRepository) resource {
	return resource{
		cache: cache,
		db:    db,
	}
}
