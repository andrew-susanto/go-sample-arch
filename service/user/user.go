package user

import (
	"github.com/joez-tkpd/go-sample-arch/entity"
)

//go:generate mockgen -source=./user.go -destination=./user_mock.go -package=user

type Service interface {
	GetUserByID(id int64) entity.User
}

type service struct {
	resource Resource
}

func NewService(rsc Resource) Service {
	return service{
		resource: rsc,
	}
}

func (svc service) GetUserByID(id int64) entity.User {
	if user := svc.resource.GetUserByIDRedis(id); user.ID > 1 {
		return user
	}

	return svc.resource.GetUserByIDPgSqlx(id)
}
