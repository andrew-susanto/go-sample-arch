package pgsqlx

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB *sqlx.DB
}

func NewSqlxDB(datasource string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", datasource)
	log.Println(err)
	return db
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		DB: db,
	}
}
