package pgsqlx

func (repo Repository) GetUserByID(id int64) User {
	return User{
		ID:     id,
		Name:   "dummy user name",
		Gender: 3,
	}
}
