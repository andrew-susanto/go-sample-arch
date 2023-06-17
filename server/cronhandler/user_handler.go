package cronhandler

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
)

// GetUserHandler gets user by given id
//
// Returns nil error if success
// Otherwise return non nil error
func (handler *Handler) GetUserHandler(ctx context.Context) error {
	_, err := handler.user.GetUser(ctx, 0)
	if err != nil {
		err = errors.Wrap(err).WithCode("HNDL.GUH00")
		log.Error(err, nil, "handler.user.GetUser() got error - GetUserHandler")
		return err
	}

	return nil
}
