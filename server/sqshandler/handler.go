package sqshandler

import (
	// go package
	"context"
	"sync"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure"
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

// InitHandler initializes sqs handler
func NewHandler(user UserUsecase) Handler {
	return Handler{
		user: user,
	}
}

// Register registers sqs handler
func (h *Handler) Register(ctx context.Context, wg *sync.WaitGroup, infra infrastructure.Infrastructure, monitor monitor.Monitor, sqs SQSService) {
	h.registerQueueConsumer(ctx, wg, monitor, infra.Config.SQSConfig.IssuanceJobFIFO, sqs, h.GetUserHandler)
}
