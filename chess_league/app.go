package chess_league

import (
	"log/slog"
	"net/http"
	"os"

	psqldb "github.com/Nag-s-Head/chess-league/db/psql_db"
	"github.com/Nag-s-Head/chess-league/handlers"
)

type App struct {
	theme Theme
	Addr  string
}

func New() *App {
	app := &App{
		theme: DefaultTheme(),
		Addr:  "0.0.0.0:8080",
	}

	return app
}

func (a *App) Run() {
	slog.Info("Starting...")
	slog.Info("Connect to the database...")
	defer os.Exit(1)

	database, err := psqldb.New()
	if err != nil {
		slog.Error("Could not connect to database - aborting", "err", err)
	}

	defer database.Close()

	err = database.GetSqlxDb().Ping()
	if err != nil {
		slog.Error("Could not ping database", "err", err)
	}

	slog.Info("Database connected successfully")
	slog.Info("Starting Nag's Knights chess league server", "addr", a.Addr)

	err = http.ListenAndServe(a.Addr, handlers.NewHandler(database))
	slog.Warn("Server has died (very sad)")
	if err != nil {
		slog.Error("Could not start", "err", err, "addr", a.Addr)
	}
}
