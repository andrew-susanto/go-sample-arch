package user

import (
	"github.com/joez-tkpd/go-sample-arch/entity"
	"github.com/joez-tkpd/go-sample-arch/repository/pgsqlx"
)

//go:generate mockgen -source=./resource_db.go -destination=./resource_db_mock.go -package=user

type DBRepository interface {
	GetUserByID(id int64) pgsqlx.User
}

func (rsc resource) GetUserByIDDB(id int64) entity.User {
	user := rsc.DB.GetUserByID(id)

	// convert to general entity object
	result := entity.User{
		ID:        user.ID,
		FirstName: user.Name,
		LastName:  "", // not provided
		Gender:    user.Gender,
	}

	return result
}
