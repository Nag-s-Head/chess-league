package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/handlers"
)

const addr = "0.0.0.0:8080"

func main() {
	slog.Info("Starting...")
	slog.Info("Connect to the database...")

	database, err := db.New()
	if err != nil {
		slog.Error("Could not connect to database - aborting", "err", err)
		os.Exit(1)
	}

	defer database.Close()

	err = database.GetSqlxDb().Ping()
	if err != nil {
		slog.Error("Could not ping database", "err", err)
		os.Exit(1)
	}

	slog.Info("Database connected successfully")
	slog.Info("Starting Nag's Knights chess league server", "addr", addr)

	err = http.ListenAndServe(addr, handlers.NewHandler())
	slog.Warn("Server has died (very sad)")
	if err != nil {
		slog.Error("Could not start", "err", err, "addr", addr)
	}
	os.Exit(1)
}
