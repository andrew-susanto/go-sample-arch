package jsonrpc

// TripItinerary is a struct that represents trip itinerary
type TripItinerary struct {
	ID                      int64                        `json:"id"`
	Time                    int64                        `json:"time"`
	Type                    string                       `json:"type"`
	Status                  string                       `json:"status"`
	PaymentStatus           string                       `json:"paymentStatus"`
	UserID                  string                       `json:"userId"`
	IssuedTime              int64                        `json:"issueTime"`
	BookingExpirationTime   int64                        `json:"bookingExpirationTime"`
	BookingContact          TripItineraryBookingContact  `json:"bookingContact"`
	AgentBookingType        string                       `json:"agentBookingType"`
	ProfileId               int64                        `json:"profileId"`
	Remarks                 *string                      `json:"remarks"`
	InvoiceID               int64                        `json:"invoiceId"`
	TotalFareWithCurrency   string                       `json:"totalFareWithCurrency"`
	Locale                  string                       `json:"locale"`
	IsCrossSelling          bool                         `json:"isCrossSelling"`
	HasGivenUpIssueAttempt  bool                         `json:"hasGivenUpIssueAttempt"`
	HasGivenUpRebookAttempt bool                         `json:"hasGivenUpRebookAttempt"`
	BookingTickets          []TripItineraryBookingTicket `json:"bookingTickets"`
}

// TripItineraryBookingContact is a struct that represents trip itinerary booking contact parts of trip itinerary
type TripItineraryBookingContact struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Phone     []string `json:"phone"`
	Email     string   `json:"email"`
}

// TripItineraryBookingTicket is a struct that represents trip itinerary booking ticket parts of trip itinerary
type TripItineraryBookingTicket struct {
	PDFTicketID string `json:"pdfTicketId"`
	Type        string `json:"type"`
}
