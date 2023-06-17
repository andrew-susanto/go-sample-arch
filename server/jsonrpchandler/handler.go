package jsonrpchandler

import (
	// go package
	"context"
	"net/http"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"github.com/andrew-susanto/go-sample-arch/usecase/account"
	"github.com/andrew-susanto/go-sample-arch/usecase/trip"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

// TripUsecase is interface for trip usecase
type TripUsecase interface {
	// GetTripItinerary gets trip itinerary by given booking id
	//
	// Returns trip itinerary and nil error if success
	// Otherwise return empty trip itinerary and non nil error
	GetTripItinerary(ctx context.Context, bookingID int64) (trip.TripItinerary, error)
}

type Handler struct {
	user UserUsecase
	trip TripUsecase
}

// InitHandler initializes json rpc handler
func NewHandler(user UserUsecase, trip TripUsecase) Handler {
	return Handler{
		user: user,
		trip: trip,
	}
}

// Register registers json rpc handler
func (h *Handler) Register(monitor monitor.Monitor) http.Handler {
	methodName := map[string]JsonRpcFunc{
		"getUserHandler":   h.GetUserHandler,
		"getTripItinerary": h.GetTripItinerary,
	}

	router := http.NewServeMux()
	router.HandleFunc("/rpc", h.handleFunc(monitor, methodName))

	handler := otelhttp.NewHandler(router, "jsonrpchandler")
	return handler
}
