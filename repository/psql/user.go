package psql

import (
	"context"

	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
)

// GetUserByID gets user by given id
//
// Return user and nil error when success
// Otherwise return empty user and non nil error
func (repo *Repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	row, err := repo.db.QueryContext(ctx, "SELECT Id FROM Users WHERE Id = $1", id)
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GUBI00")
		log.Error(err, id, "repo.db.Exec() got error - GetUserById()")
		return User{}, err
	}

	var user User
	for row.Next() {
		err = row.Scan(&user.ID)
		if err != nil {
			err = errors.Wrap(err).WithCode("RPST.GUBI01")
			log.Error(err, id, "row.Scan() got error - GetUserById()")
			return User{}, err
		}
	}

	return user, nil
}
