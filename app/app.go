package chess_league

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Nag-s-Head/chess-league/app/theme"
	psqldb "github.com/Nag-s-Head/chess-league/db/psql_db"
	"github.com/Nag-s-Head/chess-league/handlers"
)

type App struct {
	Theme theme.Theme
	Addr  string
}

func New() *App {
	app := &App{
		Theme: theme.DefaultTheme(),
		Addr:  "0.0.0.0:8080",
	}

	return app
}

func (a *App) Run() {
	slog.Info("Starting...")

	defer os.Exit(1)

	slog.Info("Connect to the database...")
	database, err := psqldb.New()
	if err != nil {
		slog.Error("Could not connect to database - aborting", "err", err)
		return
	}

	defer database.Close()

	err = database.GetSqlxDb().Ping()
	if err != nil {
		slog.Error("Could not ping database", "err", err)
		return
	}

	slog.Info("Database connected successfully")
	slog.Info("Starting Nag's Knights chess league server", "addr", a.Addr)

	server, err := handlers.NewHandler(database, a.Theme)
	if err != nil {
		slog.Error("Cannot create server handlers", "err", err, "addr", a.Addr)
	}

	err = http.ListenAndServe(a.Addr, server)
	if err != nil {
		slog.Error("Could not start", "err", err, "addr", a.Addr)
	}

	slog.Warn("Server has died (very sad)")
}
