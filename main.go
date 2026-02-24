package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const addr = "0.0.0.0:8080"
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

func main() {
	slog.Info("Starting...")
	slog.Info("Trying to connect to the database")

	var db *sqlx.DB
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
			db = pgDb
			break
		}
	}

	if db == nil {
		slog.Error("Could not connect to database - aborting")
		os.Exit(1)
	}

	slog.Info("Database connected successfully")

	slog.Info("Starting Nag's Knights chess league server", "addr", addr)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("alive and well at %s", time.Now().UTC()))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Placeholder for Nag's Knight Chess League")
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("Could not start", "err", err, "addr", addr)
	}
}
