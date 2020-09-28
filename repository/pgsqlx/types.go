package pgsqlx

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Tx = sqlx.Tx

type User struct {
	ID     int64  `db:"id"`
	Name   string `db:"full_name"`
	Gender int    `db:"gender_id"`

	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_id"`
	Status    int       `db:"status"`
}
