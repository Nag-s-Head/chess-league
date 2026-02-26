package handlers

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	privacypolicy "github.com/Nag-s-Head/chess-league/handlers/privacy_policy"
	submitgame "github.com/Nag-s-Head/chess-league/handlers/submit_game"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed index.html layout.html
var f embed.FS
var indexTmpl *template.Template = utils.GetTemplate(f, "index.html")
var layoutTmpl *template.Template = utils.GetTemplate(f, "layout.html")

type Layout struct {
	Body template.HTML
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

func Index(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := indexTmpl.Execute(&buf, nil)
	if err != nil {
		slog.Error("Cannot execute index template", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	Render(w, template.HTML(buf.String()))
}

func PrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	body, err := privacypolicy.Render()
	if err != nil {
		slog.Error("Cannot render privacy policy", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	Render(w, body)
}

func SubmitGame(w http.ResponseWriter, r *http.Request) {
	body, err := submitgame.Render()
	if err != nil {
		slog.Error("Cannot render submit game", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	Render(w, body)
}

func Test(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("alive and well at %s", time.Now().UTC())
	Render(w, template.HTML(fmt.Sprintf("<p>%s</p>", msg)))
}

// NewHandler returns a router that handles all site routes.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	// {$} matches exactly "/"
	mux.HandleFunc("GET /{$}", Index)
	mux.HandleFunc("GET /privacy-policy", PrivacyPolicy)
	mux.HandleFunc("GET /submit-game", SubmitGame)

	mux.HandleFunc("GET /test", Test)
	return mux
}
