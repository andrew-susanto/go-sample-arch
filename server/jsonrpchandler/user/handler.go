package user

import (
	"context"

	"github.com/andrew-susanto/go-sample-arch/usecase/account"
)

//go:generate mockgen -source=./handler.go -destination=./handler_mock.go -package=httphandler

// UserUsecase is interface for user usecase
type UserUsecase interface {
	// GetUser gets user by given id
	//
	// Returns user and nil error if success
	// Otherwise return empty user and non nil error
	GetUser(ctx context.Context, ID int64) (account.User, error)
}

type Handler struct {
	user UserUsecase
}

// InitHandler initializes json rpc handler
func NewHandler(user UserUsecase) Handler {
	return Handler{
		user: user,
	}
}
