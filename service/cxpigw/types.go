package cxpigw

// TripItinerary is a struct that represents trip itinerary
type TripItinerary struct {
	ID                      int64
	Time                    int64
	Type                    string
	Status                  string
	PaymentStatus           string
	UserID                  string
	IssuedTime              int64
	BookingExpirationTime   int64
	BookingContact          TripItineraryBookingContact
	AgentBookingType        string
	ProfileId               int64
	Remarks                 *string
	InvoiceID               int64
	TotalFareWithCurrency   string
	Locale                  string
	IsCrossSelling          bool
	HasGivenUpIssueAttempt  bool
	HasGivenUpRebookAttempt bool
	BookingTickets          []TripItineraryBookingTicket
}

// TripItineraryBookingContact is a struct that represents trip itinerary booking contact
type TripItineraryBookingContact struct {
	FirstName string
	LastName  string
	Phone     []string
	Email     string
}

// TripItineraryBookingTicket is a struct that represents trip itinerary booking ticket
type TripItineraryBookingTicket struct {
	PDFTicketID string
	Type        string
}
