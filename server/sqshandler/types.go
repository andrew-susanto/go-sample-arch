package sqshandler

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
