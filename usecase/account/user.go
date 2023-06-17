package account

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
)

// GetUser gets user by id
//
// Returns entity user and nil error if success
// Otherwise returns empty entity user and non-nil error
func (usecase *Usecase) GetUser(ctx context.Context, ID int64) (User, error) {
	user, err := usecase.user.GetUserByID(ctx, ID)
	if err != nil {
		err = errors.Wrap(err).WithCode("UC.GC00")
		log.Error(err, ID, "usecase.user.GetUserByID() got error - GetUser")
		return User{}, err
	}

	return User(user), nil
}
