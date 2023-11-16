package trip

import (
	"context"

	"github.com/andrew-susanto/go-sample-arch/usecase/trip"
)

//go:generate mockgen -source=./handler.go -destination=./handler_mock.go -package=httphandler

// TripUsecase is interface for trip usecase
type TripUsecase interface {
	// GetTripItinerary gets trip itinerary by given booking id
	//
	// Returns trip itinerary and nil error if success
	// Otherwise return empty trip itinerary and non nil error
	GetTripItinerary(ctx context.Context, bookingID int64) (trip.TripItinerary, error)
}

type Handler struct {
	trip TripUsecase
}

// InitHandler initializes json rpc handler
func NewHandler(trip TripUsecase) Handler {
	return Handler{
		trip: trip,
	}
}
