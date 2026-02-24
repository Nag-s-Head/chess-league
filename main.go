package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
)

const addr = "0.0.0.0:8080"

func main() {
	slog.Info("Starting...")
	slog.Info("Connect to the database...")

	database, err := db.Connect()
	if err != nil {
		slog.Error("Could not connect to database - aborting", "err", err)
		os.Exit(1)
	}

	err = database.Ping()
	if err != nil {
		slog.Error("Could not ping database", "err", err)
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

	err = http.ListenAndServe(addr, nil)
	slog.Warn("Server has died (very sad)")
	if err != nil {
		slog.Error("Could not start", "err", err, "addr", addr)
	}
	os.Exit(1)
}
