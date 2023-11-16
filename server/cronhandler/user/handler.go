package user

import (
	"context"

	"github.com/andrew-susanto/go-sample-arch/usecase/account"
)

//go:generate mockgen -source=./handler.go -destination=./handler_mock.go -package=httphandler
type UserUsecase interface {
	GetUser(ctx context.Context, ID int64) (account.User, error)
}

type Handler struct {
	user UserUsecase
}

func NewHandler(user UserUsecase) Handler {
	return Handler{
		user: user,
	}
}
