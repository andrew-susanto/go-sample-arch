package sqshandler

import (
	// go package
	"context"
	"sync"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

//go:generate mockgen -source=./handler.go -destination=./handler_mock.go -package=httphandler
type UserHandler interface {
	GetUserHandler(ctx context.Context, message types.Message) error
}

type Handler struct {
	user UserHandler
}

// InitHandler initializes sqs handler
func NewHandler(user UserHandler) Handler {
	return Handler{
		user: user,
	}
}

// Register registers sqs handler
func (h *Handler) Register(ctx context.Context, wg *sync.WaitGroup, infra infrastructure.Infrastructure, monitor monitor.Monitor, sqs SQSService) {
	h.registerQueueConsumer(ctx, wg, monitor, infra.Config.SQSConfig.IssuanceJobFIFO, sqs, h.user.GetUserHandler)
}
