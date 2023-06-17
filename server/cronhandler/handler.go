package cronhandler

import (
	// go package
	"context"
	"sync"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
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

// Register registers cron handler
func (h *Handler) Register(ctx context.Context, wg *sync.WaitGroup, monitor monitor.Monitor) {
	h.registerCron(ctx, wg, monitor, "testCronGetUserHandlerEvenryHour", 3600, h.GetUserHandler)
}
