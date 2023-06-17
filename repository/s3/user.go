package s3

import (
	// golang package
	"context"
)

// GetUserByID is a function to get user by id
//
// Returns User and nil error if success
// Otherwise return empty User and non nil error
func (repo *Repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	// repo.download.DownloadWithContext(ctx, ...)

	return User{
		ID:     id,
		Name:   "dummy user name",
		Gender: 3,
	}, nil
}
