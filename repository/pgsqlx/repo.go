package pgsqlx

import (
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=pgsqlx

type Repository interface {
	GetUserByID(id int64) User
}

type repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{
		DB: db,
	}
}

func (repo repository) GetUserByID(id int64) User {
	return User{
		ID:     0,
		Name:   "",
		Gender: 0,
	}
}
