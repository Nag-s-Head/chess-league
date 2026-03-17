package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	playerdetails "github.com/Nag-s-Head/chess-league/handlers/player_details"
	submitgame "github.com/Nag-s-Head/chess-league/handlers/submit_game"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed index.html layout.html
var f embed.FS
var indexTmpl *template.Template = utils.GetTemplate(f, "index.html")
var layoutTmpl *template.Template = utils.GetTemplate(f, "layout.html")

type Layout struct {
	Body template.HTML
}

type IndexData struct {
	Players      []model.Player
	TotalGames   int
	TotalPlayers int
}

// Render wraps the provided body HTML in the global layout and writes it to w.
func Render(w http.ResponseWriter, body template.HTML) {
	err := layoutTmpl.Execute(w, Layout{
		Body: body,
	})
	if err != nil {
		slog.Error("Cannot execute layout template", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("alive and well at %s", time.Now().UTC())
	Render(w, template.HTML(fmt.Sprintf("<p>%s</p>", msg)))
}

func PlayerDetails(db *db.Db) func(w http.ResponseWriter, r *http.Request) {
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

		Render(w, body)
	}
}

// NewHandler returns a router that handles all site routes.
func NewHandler(db *db.Db) http.Handler {
	mux := http.NewServeMux()
	// {$} matches exactly "/"
	mux.HandleFunc("GET /{$}", Index(db))
	mux.HandleFunc("GET /player/{id}", PlayerDetails(db))
	mux.HandleFunc("GET /test", Test)
	mux.HandleFunc("GET /privacy-policy", PrivacyPolicy)
	mux.HandleFunc("GET /rules", Rules)
	mux.HandleFunc("GET /rules/agree", RulesAgree)

	slog.Info(fmt.Sprintf("To submit a game use http://0.0.0.0:8080/%s?%s=%s",
		submitgame.BasePath,
		submitgame.MagicNumberParam,
		magicNumber))
	mux.HandleFunc(fmt.Sprintf("GET %s", submitgame.BasePath), SubmitGame(db))
	submitgame.Register(mux, db)

	return mux
}
