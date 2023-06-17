package http

import "context"

// GetUserByID gets user by given id from http
//
// Returns user and nil error when success
// Otherwise return empty user and non nil error
func (repository *Repository) GetUserByID(ctx context.Context, userID int64) (User, error) {
	// repository.httpClient.Do()

	user := User{
		ID: userID,
	}
	return user, nil
}
