package chess_league

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	psqldb "github.com/Nag-s-Head/chess-league/db/psql_db"
	"github.com/Nag-s-Head/chess-league/handlers"
)

//go:embed theme.css
var fs embed.FS

type App struct {
	Theme Theme
	Addr  string
}

func New() *App {
	app := &App{
		Theme: DefaultTheme(),
		Addr:  "0.0.0.0:8080",
	}

	return app
}

func (a *App) generateThemeCss() ([]byte, error) {
	theme, err := fs.ReadFile("theme.css")
	if err != nil {
		return nil, errors.Join(errors.New("Cannot read theme.css"), err)
	}

	tmpl, err := template.New("css").Parse(string(theme))
	if err != nil {
		return nil, errors.Join(errors.New("Cannot parse theme.css"), err)
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, a.Theme)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute theme.css template"), err)
	}

	return buf.Bytes(), nil
}

func (a *App) Run() {
	slog.Info("Starting...")

	defer os.Exit(1)

	slog.Info("Generating theme css...")
	themeCss, err := a.generateThemeCss()
	if err != nil {
		slog.Error("Could not generate theme CSS", "error", err)
		return
	}

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

	err = http.ListenAndServe(a.Addr, handlers.NewHandler(database, themeCss))
	if err != nil {
		slog.Error("Could not start", "err", err, "addr", a.Addr)
	}

	slog.Warn("Server has died (very sad)")
}
