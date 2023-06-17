package trip

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/service/cxpigw"
)

//go:generate mockgen -source=./usecase.go -destination=./usecase_mock.go -package=account

type CxpIgwService interface {
	// GetTripItineraryByBookingID gets trip itinerary by booking id
	//
	// Returns TripItinerary and nil error if sucess
	// Otherwise return empty TripItinerary and non nil error
	GetTripItineraryByBookingID(ctx context.Context, bookingID int64) (cxpigw.TripItinerary, error)
}

type Usecase struct {
	cxpigw CxpIgwService
}

// NewUsecase creates new trip usecase
func NewUsecase(cxpigw CxpIgwService) Usecase {
	return Usecase{
		cxpigw: cxpigw,
	}
}
