package cache

// User is a struct that contains user data
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    int    `json:"gender"`
}
