package jsonrpchandler

import (
	// go package

	"encoding/json"
	"net/http"

	// internal package

	"github.com/andrew-susanto/go-sample-arch/infrastructure/jsonrpccontext"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

//go:generate mockgen -source=./handler.go -destination=./handler_mock.go -package=httphandler

// UserUsecase is interface for user usecase
type UserHandler interface {
	// GetUser gets user by given id
	//
	// Returns user and nil error if success
	// Otherwise return empty user and non nil error
	GetUserHandler(tdkCtx *jsonrpccontext.TdkJsonRpcContext, params json.RawMessage) (interface{}, error)
}

type TripHandler interface {
	GetTripItinerary(tdkCtx *jsonrpccontext.TdkJsonRpcContext, params json.RawMessage) (interface{}, error)
}

type Handler struct {
	user    UserHandler
	trip    TripHandler
	monitor monitor.Monitor
}

// InitHandler initializes json rpc handler
func NewHandler(user UserHandler, trip TripHandler, monitor monitor.Monitor) Handler {
	return Handler{
		user:    user,
		trip:    trip,
		monitor: monitor,
	}
}

// Register registers json rpc handler
func (h *Handler) Register() http.Handler {
	methodName := map[string]JsonRpcFunc{
		"getUserHandler":   h.user.GetUserHandler,
		"getTripItinerary": h.trip.GetTripItinerary,
	}

	router := http.NewServeMux()
	router.HandleFunc("/rpc", h.handleFunc(methodName))

	handler := otelhttp.NewHandler(router, "jsonrpchandler")
	return handler
}
