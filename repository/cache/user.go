package cache

import (
	// golang package
	"context"
)

const (
	keyUser = "user_"
)

// GetUserByID gets user by given id from cache
//
// Returns user and nil error when success
// Otherwise return empty user and non nil error
func (repo *Repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	// resp := repo.redis.Get(ctx, "123")

	user := User{}
	return user, nil
}
