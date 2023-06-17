package httphandler

import (
	// go package
	"context"
	"net/http"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/httpcontext"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"github.com/andrew-susanto/go-sample-arch/usecase/account"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

//go:generate mockgen -source=./handler.go -destination=./handler_mock.go -package=httphandler
type UserUsecase interface {
	GetUser(ctx context.Context, ID int64) (account.User, error)
}

type Handler struct {
	user UserUsecase
}

// NewHandler initializes new http handler based on given usecase
func NewHandler(user UserUsecase) Handler {
	return Handler{
		user: user,
	}
}

// Register registers http handler
func (h *Handler) Register(monitor monitor.Monitor) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/health", h.handleFunc(monitor, h.HealthCheckHandler))
	router.HandleFunc("/user/{id}", h.handleFunc(monitor, h.GetUserHandler))

	handler := otelhttp.NewHandler(router, "httphandler")
	return handler
}

// HealthCheckHandler handles health check request
func (h *Handler) HealthCheckHandler(tdkCtx *httpcontext.TdkHttpContext) error {
	tdkCtx.WriteHTTPResponseToJSON(map[string]interface{}{
		"status": "OK",
	}, http.StatusOK)
	return nil
}
