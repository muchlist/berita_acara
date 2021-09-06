package dao

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/berita_acara/db"
)

type TxOpt struct {
	Ctx context.Context
	Tx  SqlGdbc
}

func GetTx(ctx context.Context) (SqlGdbc, error) {
	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SqlGdbc (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It should be able to work with all SQL data that follows SQL standard.
type SqlGdbc interface {
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Transactioner
}

// Transactioner is the transaction interface for database handler
// It should only be applicable to SQL database
type Transactioner interface {
	// Rollback a transaction
	Rollback(ctx context.Context) error
	// Commit a transaction
	Commit(ctx context.Context) error
	// Begin gets *sql.DB from receiver and return a SqlGdbc, which has a *sql.Tx
	Begin(ctx context.Context) (pgx.Tx, error)
}

// SqlDBTx is the concrete implementation of sqlGdbc by using *sql.DB
type SqlDBTx struct {
	DB *pgxpool.Pool
}

// SqlConnTx is the concrete implementation of sqlGdbc by using *sql.Tx
type SqlConnTx struct {
	DB *pgx.Tx
}
