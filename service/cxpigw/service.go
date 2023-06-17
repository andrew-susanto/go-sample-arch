package cxpigw

import "context"

//go:generate mockgen -source=./service.go -destination=./service_mock.go -package=user

type Resource interface {
	// GetTripItineraryByBookingIDFromJSONRpc is a function to get trip itinerary by booking id from jsonrpc
	//
	// Returns TripItinerary and nil error if sucess
	// Otherwise return empty TripItinerary and non nil error
	GetTripItineraryByBookingIDFromJSONRpc(ctx context.Context, bookingID int64) (TripItinerary, error)
}

type Service struct {
	resource Resource
}

// NewService creates new service for cxpigw service
func NewService(rsc Resource) Service {
	return Service{
		resource: rsc,
	}
}
