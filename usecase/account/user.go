package account

import (
	"github.com/joez-tkpd/go-sample-arch/entity"
)

//go:generate mockgen -source=./user.go -destination=./user_mock.go -package=account

type UserService interface {
	GetUserByID(id int64) entity.User
}

func (uc Usecase) GetUser(id int64) entity.User {
	return uc.user.GetUserByID(id)
}
