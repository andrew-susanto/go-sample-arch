package sns

import (
	"time"
)

// User is a struct to represent user data
type User struct {
	ID     int64  `json:"id"`
	Name   string `json:"full_name"`
	Gender int    `json:"gender_id"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_id"`
	Status    int       `json:"status"`
}
