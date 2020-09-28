package account

import (
	"github.com/joez-tkpd/go-sample-arch/entity"
	"github.com/joez-tkpd/go-sample-arch/service/user"
)

//go:generate mockgen -source=./user.go -destination=./user_mock.go -package=account

type Usecase interface {
	GetUser(id int64) entity.User
}

type usecase struct {
	user user.Service
}

func NewUsecase(user user.Service) Usecase {
	return usecase{
		user: user,
	}
}

func (uc usecase) GetUser(id int64) entity.User {
	return uc.user.GetUserByID(id)
}
