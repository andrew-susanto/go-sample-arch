package account

type Usecase struct {
	user UserService
}

func NewUsecase(user UserService) Usecase {
	return Usecase{
		user: user,
	}
}
