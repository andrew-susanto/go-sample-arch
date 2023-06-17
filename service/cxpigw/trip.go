package cxpigw

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
)

// GetTripItineraryByBookingID is a function to get trip itinerary by booking id
//
// Returns TripItinerary and nil error if sucess
// Otherwise return empty TripItinerary and non nil error
func (svc *Service) GetTripItineraryByBookingID(ctx context.Context, bookingID int64) (TripItinerary, error) {
	ctx, span := tracer.Start(ctx, "service.cxpigw.GetTripItineraryByBookingID")
	defer span.End()

	resp, err := svc.resource.GetTripItineraryByBookingIDFromJSONRpc(ctx, bookingID)
	if err != nil {
		err = errors.Wrap(err).WithCode("SVC.GTIBBID00")
		log.Error(err, nil, "svc.resource.GetTripItineraryByBookingIDFromJSONRpc() got error - GetTripItineraryByBookingID")
		return TripItinerary{}, err
	}

	return resp, nil
}
