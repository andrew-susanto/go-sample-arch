package trip

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
)

// GetTripItinerary gets trip itinerary by booking id
//
// Returns TripItinerary and nil error if sucess
// Otherwise return empty TripItinerary and non nil error
func (usecase *Usecase) GetTripItinerary(ctx context.Context, bookingID int64) (TripItinerary, error) {
	ctx, span := tracer.Start(ctx, "usecase.trip.GetTripItinerary")
	defer span.End()

	resp, err := usecase.cxpigw.GetTripItineraryByBookingID(ctx, bookingID)
	if err != nil {
		err = errors.Wrap(err).WithCode("UC.GTT00")
		log.Error(err, bookingID, "usecase.cxpigw.GetTripItineraryByBookingID() got error - GetTripItinerary")
		return TripItinerary{}, err
	}

	tripItineraryBookingContact := TripItineraryBookingContact(resp.BookingContact)
	tripItinearyBookingTickets := make([]TripItineraryBookingTicket, len(resp.BookingTickets))
	for i := 0; i < len(resp.BookingTickets); i++ {
		tripItinearyBookingTickets[i] = TripItineraryBookingTicket(resp.BookingTickets[i])
	}

	tripItinerary := TripItinerary{
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
