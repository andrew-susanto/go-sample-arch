package db

import (
	// golang package
	"context"
	"database/sql"
	"fmt"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/secretmanager"

	// external package
	"github.com/jmoiron/sqlx"
)

// Database is interface of database connection
type Database interface {
	// Exec executes a query without returning any rows. The args are for any placeholder parameters in the query.
	//
	// Exec uses context.Background internally; to specify the context, use ExecContext.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// Query executes a query that returns rows, typically a SELECT. The args are for any placeholder parameters in the query.
	//
	// Query uses context.Background internally; to specify the context, use QueryContext.
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow executes a query that is expected to return at most one row. QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
	//
	// QueryRow uses context.Background internally; to specify the context, use QueryRowContext.
	QueryRow(query string, args ...interface{}) *sql.Row

	// QueryContext executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// BeginTx starts a transaction.
	//
	// The provided context is used until the transaction is committed or rolled back. If the context is canceled, the sql package will roll back the transaction. Tx.Commit will return an error if the context provided to BeginTx is canceled.
	//
	// The provided TxOptions is optional and may be nil if defaults should be used. If a non-default isolation level is used that the driver doesn't support, an error will be returned.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// Close closes the database and prevents new queries from starting.
	// Close then waits for all queries that have started processing on the server
	// to finish.
	//
	// It is rare to Close a DB, as the DB handle is meant to be
	// long-lived and shared between many goroutines.
	Close() error
}

// InitDatabaseClient initialize database client
func InitDatabaseClient(config secretmanager.SecretsPostgreSQL) Database {
	var err error

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.DBName)

	// connect database
	psqlDb, err := sqlx.Connect("postgres", psqlconn)
	if err != nil {
		log.Fatal(err, nil, "sqlx.Connect() got error - InitDatabaseClient")
		return psqlDb
	}

	// check db
	err = psqlDb.Ping()
	if err != nil {
		log.Fatal(err, nil, "psqlDb.Ping() got error - InitDatabaseClient")
		return psqlDb
	}

	log.Info(nil, "success connect to postgres database - InitDatabaseClient")
	return psqlDb
}

// CloseDBConnection close the database connection
func CloseDBConnection(db Database) {
	err := db.Close()
	if err != nil {
		log.Error(err, nil, "db.Close() got error - CloseDBConnection")
	}

	log.Info(nil, "Success close db connection")
}
