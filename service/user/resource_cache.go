package user

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/entity"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
)

// GetUserByIDFromCache gets user by id from cache service
//
// Returns entity user and nil error if success
// Otherwise returns empty entity user and non-nil error
func (resource *resource) GetUserByIDFromCache(ctx context.Context, ID int64) (entity.User, error) {
	user, err := resource.cache.GetUserByID(ctx, ID)
	if err != nil {
		err = errors.Wrap(err).WithCode("RSC.GUBIFC00")
		log.Error(err, ID, "resource.cache.GetUserByID() got erro - GetUserByIDFromCache")
		return entity.User{}, err
	}

	// convert to general entity object
	result := entity.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    user.Gender,
	}

	return result, nil
}
