package user

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/entity"
)

//go:generate mockgen -source=./service.go -destination=./service_mock.go -package=user

type Resource interface {
	// GetUserByIDFromDB gets user by id from db service
	//
	// Returns entity user and nil error if success
	// Otherwise returns empty entity user and non-nil error
	GetUserByIDFromDB(ctx context.Context, ID int64) (entity.User, error)

	// GetUserByIDFromCache gets user by id from cache service
	//
	// Returns entity user and nil error if success
	// Otherwise returns empty entity user and non-nil error
	GetUserByIDFromCache(ctx context.Context, ID int64) (entity.User, error)
}

type Service struct {
	resource Resource
}

// NewService creates new user service
func NewService(rsc Resource) Service {
	return Service{
		resource: rsc,
	}
}
