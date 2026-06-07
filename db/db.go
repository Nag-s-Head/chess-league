package db

import "github.com/jmoiron/sqlx"

type Db interface {
	GetSqlxDb() *sqlx.DB
	Close()
	DoTx(fn func(tx *sqlx.Tx) error) error
}
