package user

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/entity"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
)

// GetUserByIDFromDB gets user by id from db service
//
// Returns entity user and nil error if success
// Otherwise returns empty entity user and non-nil error
func (resource *resource) GetUserByIDFromDB(ctx context.Context, ID int64) (entity.User, error) {
	user, err := resource.db.GetUserByID(ctx, ID)
	if err != nil {
		err = errors.Wrap(err).WithCode("RSC.GUBIFDB00")
		log.Error(err, nil, "resource.db.GetUserByID() got error - GetUserByIDFromDB")
		return entity.User{}, err
	}

	// convert to general entity object
	result := entity.User{
		ID:        user.ID,
		FirstName: user.Name,
		LastName:  "",
		Gender:    user.Gender,
	}

	return result, nil
}
