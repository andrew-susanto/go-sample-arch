package jsonrpc

import (
	// goland package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
)

const (
	getTripItineraryMethodName = "getTripItinerary"
)

// GetTripItinerary gets trip itinerary by booking id
//
// Returns TripItinerary and nil error if sucess
// Otherwise return empty TripItinerary and non nil error
func (repo *Repository) GetTripItinerary(ctx context.Context, bookingID int64) (TripItinerary, error) {
	ctx, span := tracer.Start(ctx, "repository.jsonrpc.GetTripItinerary")
	defer span.End()

	resp, err := repo.cxpigw.Call(ctx, getTripItineraryMethodName, bookingID)
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GTT00")
		log.Error(err, bookingID, "repo.cxpigw.Call() got error - GetTripItinerary")
		return TripItinerary{}, err
	}

	var tripItinerary TripItinerary
	err = resp.GetObject(&tripItinerary)
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GTT01")
		log.Error(err, bookingID, "resp.GetObject() got error - GetTripItinerary")
		return TripItinerary{}, err
	}

	return tripItinerary, nil
}
