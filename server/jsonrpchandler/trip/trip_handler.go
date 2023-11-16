package trip

import (
	// golang package
	"encoding/json"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/jsonrpccontext"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
)

// GetTripItinerary gets trip itinerary by booking id
//
// Returns TripItinerary and nil error if sucess
// Otherwise return empty TripItinerary and non nil error
func (handler *Handler) GetTripItinerary(tdkCtx *jsonrpccontext.TdkJsonRpcContext, params json.RawMessage) (interface{}, error) {
	ctx, span := tracer.Start(tdkCtx.Context, "handler.jsonrpchandler.GetTripItinerary")
	defer span.End()

	var bookingID []int64

	err := json.Unmarshal(params, &bookingID)
	if err != nil {
		err = errors.Wrap(err).WithCode("HNDL.GTT00").WithType(errors.USER)
		log.Error(err, nil, "json.Umarhsal() got error - GetTripItinerary")
		return nil, err
	}

	if len(bookingID) == 0 {
		err = errors.New("invalid parameter").WithCode("HNDL.GTT01").WithType(errors.USER)
		log.Error(err, nil, "invalid parameter - GetTripItinerary")
		return nil, err
	}

	resp, err := handler.trip.GetTripItinerary(ctx, bookingID[0])
	if err != nil {
		err = errors.Wrap(err).WithCode("HNDL.GUH00")
		log.Error(err, nil, "handler.trip.GetTripItinerary() got error - GetTripItinerary")
		return nil, err
	}

	tripItineraryBookingContact := TripItineraryBookingContact(resp.BookingContact)
	tripItinearyBookingTickets := make([]TripItineraryBookingTicket, len(resp.BookingTickets))
	for i := 0; i < len(resp.BookingTickets); i++ {
		tripItinearyBookingTickets[i] = TripItineraryBookingTicket(resp.BookingTickets[i])
	}

	tripItinerary := TripItineraryResponse{
		ID:                      resp.ID,
		Time:                    resp.Time,
		Type:                    resp.Type,
		Status:                  resp.Status,
		PaymentStatus:           resp.PaymentStatus,
		UserID:                  resp.UserID,
		IssuedTime:              resp.IssuedTime,
		BookingExpirationTime:   resp.BookingExpirationTime,
		BookingContact:          tripItineraryBookingContact,
		AgentBookingType:        resp.AgentBookingType,
		ProfileId:               resp.ProfileId,
		Remarks:                 resp.Remarks,
		InvoiceID:               resp.InvoiceID,
		TotalFareWithCurrency:   resp.TotalFareWithCurrency,
		Locale:                  resp.Locale,
		IsCrossSelling:          resp.IsCrossSelling,
		HasGivenUpIssueAttempt:  resp.HasGivenUpIssueAttempt,
		HasGivenUpRebookAttempt: resp.HasGivenUpRebookAttempt,
		BookingTickets:          tripItinearyBookingTickets,
	}

	return tripItinerary, nil
}
