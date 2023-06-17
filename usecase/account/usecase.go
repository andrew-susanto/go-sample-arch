package account

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/service/user"
)

//go:generate mockgen -source=./usecase.go -destination=./usecase_mock.go -package=account

type UserService interface {
	// GetUserByID gets user by id
	//
	// Returns entity user and nil error if success
	// Otherwise returns empty entity user and non-nil error
	GetUserByID(ctx context.Context, ID int64) (user.User, error)
}

type Usecase struct {
	user UserService
}

// NewUsecase creates new account usecase
func NewUsecase(user UserService) Usecase {
	return Usecase{
		user: user,
	}
}
