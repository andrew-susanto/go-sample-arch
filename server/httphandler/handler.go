package httphandler

import (
	// go package

	"net/http"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/httpcontext"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

//go:generate mockgen -source=./handler.go -destination=./handler_mock.go -package=httphandler
type UserHandler interface {
	// HandleGetUser gets user by given id
	//
	// Returns user and nil error if success
	// Otherwise return empty user and non nil error
	HandleGetUser(tdkCtx *httpcontext.TdkHttpContext) error
}

type Handler struct {
	user    UserHandler
	monitor monitor.Monitor
}

// NewHandler initializes new http handler based on given usecase
func NewHandler(user UserHandler, monitor monitor.Monitor) Handler {
	return Handler{
		user:    user,
		monitor: monitor,
	}
}

// Register registers http handler
func (handler *Handler) Register() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/health", handler.handleFunc(handler.HealthCheckHandler))
	router.HandleFunc("/user/{id}", handler.handleFunc(handler.user.HandleGetUser))

	otelHandler := otelhttp.NewHandler(router, "httphandler")
	return otelHandler
}

// HealthCheckHandler handles health check request
func (h *Handler) HealthCheckHandler(tdkCtx *httpcontext.TdkHttpContext) error {
	tdkCtx.WriteHTTPResponseToJSON(map[string]interface{}{
		"status": "OK",
	}, http.StatusOK)
	return nil
}
