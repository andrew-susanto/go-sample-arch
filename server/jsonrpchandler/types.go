package jsonrpchandler

// GetUserParam is a struct that represents request parameter for GetUserHandler
type GetUserParam struct {
	ID int64 `json:"id"`
}

// GetUserResponse is a struct that represents response for GetUserHandler
type GetUserResponse struct {
	FirstName string `json:"firstName"`
	Gender    int    `json:"gender"`
	ID        int64  `json:"id"`
	LastName  string `json:"lastName"`
}

// GetTripItineraryParam is a struct that represents request parameter for GetTripItinerary
type TripItineraryResponse struct {
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

// TripItinerary is a struct that represents response for GetTripItinerary
type TripItineraryBookingContact struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Phone     []string `json:"phone"`
	Email     string   `json:"email"`
}

// TripItineraryBookingTicket is a struct that represents response for GetTripItinerary
type TripItineraryBookingTicket struct {
	PDFTicketID string `json:"pdfTicketId"`
	Type        string `json:"type"`
}
