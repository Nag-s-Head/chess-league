package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/admin"
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
	Body    template.HTML
	IsAdmin bool
}

type IndexData struct {
	Players      []model.Player
	TotalGames   int
	TotalPlayers int
}

// WithLayout wraps the provided body HTML in the global layout and writes it to w.
func withLayout(w http.ResponseWriter, body template.HTML, isAdmin bool) {
	err := layoutTmpl.Execute(w, Layout{
		Body:    body,
		IsAdmin: isAdmin,
	})
	if err != nil {
		slog.Error("Cannot execute layout template", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func WithLayoutAdmin(w http.ResponseWriter, body template.HTML) {
	withLayout(w, body, true)
}

func WithLayout(w http.ResponseWriter, body template.HTML) {
	withLayout(w, body, false)
}

func Test(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("alive and well at %s", time.Now().UTC())
	WithLayout(w, template.HTML(fmt.Sprintf("<p>%s</p>", msg)))
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

		WithLayout(w, body)
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
	mux.HandleFunc(fmt.Sprintf("GET %s", submitgame.BasePath), SubmitGame(db))
	submitgame.Register(mux, db)
	admin.Register(mux, db, WithLayoutAdmin)

	slog.Info(fmt.Sprintf("To submit a game use %s/%s?%s=%s",
		os.Getenv("APP_BASE_URL"),
		submitgame.BasePath,
		submitgame.MagicNumberParam,
		magicNumber))

	return mux
}
