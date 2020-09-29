package user

import (
	"github.com/joez-tkpd/go-sample-arch/entity"
)

type Service struct {
	resource Resource
}

func NewService(rsc Resource) Service {
	return Service{
		resource: rsc,
	}
}

func (svc Service) GetUserByID(id int64) entity.User {
	if user := svc.resource.GetUserByIDCache(id); user.ID > 1 {
		return user
	}

	return svc.resource.GetUserByIDDB(id)
}
