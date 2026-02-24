package db

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

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
		slog.Info("Trying to connect to database...", "try number", tries+1, "max tries", maxTries)
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
