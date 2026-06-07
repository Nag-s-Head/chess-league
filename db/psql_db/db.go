package psqldb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Nag-s-Head/chess-league/db/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const databaseEnvVarName = "DATABASE_URL"

func connect() (*sqlx.DB, error) {
	url := os.Getenv(databaseEnvVarName)
	if url == "" {
		return nil, fmt.Errorf("Env var %s is not set", databaseEnvVarName)
	}

	conn, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Cannot connect to database"), err)
	}
	return conn, nil
}

func InternalConnect() (*sqlx.DB, error) {
	var database *sqlx.DB
	const (
		maxTries = 10
		wait     = time.Second / 2
	)

	for tries := range maxTries {
		tries++
		slog.Info("Trying to connect to database...", "try number", tries, "max tries", maxTries)
		pgDb, err := connect()
		if err != nil {
			slog.Warn("Could not connect to database trying again...", "waiting for", wait, "err", err)
			time.Sleep(wait)
		} else {
			database = pgDb
			break
		}
	}

	if database == nil {
		slog.Error("Could not connect to database - aborting")
		return nil, errors.New("Cannot connect")
	}
	return database, nil
}

type Db struct {
	db *sqlx.DB
}

func New() (*Db, error) {
	db, err := InternalConnect()
	if err != nil {
		return nil, err
	}

	return From(db)
}

func From(in *sqlx.DB) (*Db, error) {
	db := &Db{db: in}
	migrator := migrations.Migrations()
	migrator.Rebinder = sqlx.DOLLAR
	err := migrator.Migrate(db)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot migrate database"), err)
	}

	return db, nil
}

func (d *Db) GetSqlxDb() *sqlx.DB {
	return d.db
}

func (d *Db) Close() {
	d.db.Close()
}

func (d *Db) DoTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := d.GetSqlxDb().BeginTxx(context.Background(), nil)
	if err != nil {
		return errors.Join(errors.New("Cannot start transaction"), err)
	}

	defer tx.Rollback()
	err = fn(tx)
	if err != nil {
		return errors.Join(errors.New("There was an error whilst executing the transaction"), err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("Cannot commit transaction"), err)
	}

	return nil
}
