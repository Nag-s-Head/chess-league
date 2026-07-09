package handlers

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Nag-s-Head/chess-league/chess_league/theme"
	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/handlers/admin"
	"github.com/Nag-s-Head/chess-league/handlers/assets"
	"github.com/Nag-s-Head/chess-league/handlers/league"
	playerdetails "github.com/Nag-s-Head/chess-league/handlers/player_details"
	submitgame "github.com/Nag-s-Head/chess-league/handlers/submit_game"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed index.html layout.html theme.css
var f embed.FS
var indexTmpl *template.Template = utils.GetTemplate(f, "index.html")
var layoutTmpl *template.Template = utils.GetTemplate(f, "layout.html")
var themeCssTmpl *template.Template = utils.GetTemplate(f, "theme.css")

type Layout struct {
	Body    template.HTML
	IsAdmin bool
	Theme   theme.Theme
}

func generateThemeCss(theme theme.Theme) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := themeCssTmpl.Execute(buf, theme)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute theme.css template"), err)
	}

	return buf.Bytes(), nil
}

// WithLayout wraps the provided body HTML in the global layout and writes it to w.
func withLayout(w http.ResponseWriter, body template.HTML, isAdmin bool, theme theme.Theme) {
	err := layoutTmpl.Execute(w, Layout{
		Body:    body,
		IsAdmin: isAdmin,
		Theme:   theme,
	})
	if err != nil {
		slog.Error("Cannot execute layout template", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type LayoutFn func(w http.ResponseWriter, body template.HTML)

func WithLayoutAdmin(theme theme.Theme) LayoutFn {
	return func(w http.ResponseWriter, body template.HTML) {
		withLayout(w, body, true, theme)
	}
}

func WithLayout(theme theme.Theme) LayoutFn {
	return func(w http.ResponseWriter, body template.HTML) {

		withLayout(w, body, false, theme)
	}
}

func Test(WithLayout LayoutFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("alive and well at %s", time.Now().UTC())
		WithLayout(w, template.HTML(fmt.Sprintf("<p>%s</p>", msg)))
	}
}

func PlayerDetails(db db.Db, WithLayout LayoutFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			slog.Warn("Invalid player ID", "id", idStr)
			http.Error(w, "Invalid player ID", http.StatusBadRequest)
			return
		}

		body, err := playerdetails.Render(db, id)
		if err != nil {
			slog.Error("Cannot render player details", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		WithLayout(w, body)
	}
}

func League(db db.Db, WithLayout LayoutFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := league.Render(db)
		if err != nil {
			slog.Error("Cannot render league page", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		WithLayout(w, body)
	}
}

// NewHandler returns a router that handles all site routes.
func NewHandler(db db.Db, theme theme.Theme) (http.Handler, error) {
	themeCss, err := generateThemeCss(theme)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot generate theme css"), err)
	}

	mux := http.NewServeMux()
	layoutFn := WithLayout(theme)
	// {$} matches exactly "/"
	mux.HandleFunc("GET /{$}", Index(db, theme))
	mux.HandleFunc("GET /player/{id}", PlayerDetails(db, layoutFn))
	mux.HandleFunc("GET /test", Test(layoutFn))
	mux.HandleFunc("GET /privacy-policy", PrivacyPolicy(layoutFn))
	mux.HandleFunc("GET /league", League(db, layoutFn))
	mux.HandleFunc("GET /rules", Rules(layoutFn))
	mux.HandleFunc("GET /rules/agree", RulesAgree)
	mux.HandleFunc(fmt.Sprintf("GET %s", submitgame.BasePath), SubmitGame(db, layoutFn))
	submitgame.Register(mux, db)
	admin.Register(mux, db, WithLayoutAdmin(theme))
	assets.Register(mux, themeCss)

	slog.Info(fmt.Sprintf("To submit a game use %s/%s?%s=%s",
		os.Getenv("APP_BASE_URL"),
		submitgame.BasePath,
		submitgame.MagicNumberParam,
		magicNumber))

	return mux, nil
}
