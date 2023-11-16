package cronhandler

import (
	// go package
	"context"
	"sync"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
)

//go:generate mockgen -source=./handler.go -destination=./handler_mock.go -package=httphandler
type UserHandler interface {
	GetUserHandler(ctx context.Context) error
}

type Handler struct {
	user UserHandler
}

func NewHandler(user UserHandler) Handler {
	return Handler{
		user: user,
	}
}

// Register registers cron handler
func (h *Handler) Register(ctx context.Context, wg *sync.WaitGroup, monitor monitor.Monitor) {
	h.registerCron(ctx, wg, monitor, "testCronGetUserHandlerEvenryHour", 3600, h.user.GetUserHandler)
}
