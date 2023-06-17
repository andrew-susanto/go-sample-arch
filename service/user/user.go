package user

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
)

// GetUserByID gets user by id
//
// Returns entity user and nil error if success
// Otherwise returns empty entity user and non-nil error
func (svc *Service) GetUserByID(ctx context.Context, ID int64) (User, error) {
	user, err := svc.resource.GetUserByIDFromCache(ctx, ID)
	if err != nil {
		err = errors.Wrap(err).WithCode("SVC.GUBI00")
		log.Error(err, nil, "svc.resource.GetUserByIDFromCache() got error - GetUserByID")
		return User{}, err
	}

	if user.ID > 0 {
		return User(user), nil
	}

	user, err = svc.resource.GetUserByIDFromDB(ctx, ID)
	if err != nil {
		err = errors.Wrap(err).WithCode("SVC.GUBI01")
		log.Error(err, nil, "svc.resource.GetUserByIDFromDB() got error - GetUserByID")
		return User{}, err
	}

	return User(user), nil
}
