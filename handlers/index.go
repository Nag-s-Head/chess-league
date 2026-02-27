package handlers

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
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

func Index(db *db.Db) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		players, err := model.GetPlayersByElo(db)
		if err != nil {
			slog.Warn("Could not get leaderboard", "err", err)
		}

		var buf bytes.Buffer
		err = indexTmpl.Execute(&buf, players)
		if err != nil {
			slog.Error("Cannot execute index template", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		Render(w, template.HTML(buf.String()))
	}
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

var magicNumber string = os.Getenv(submitgame.MagicNumberEnvVar)

// probably long enough to submit a game
const maxAge = 3 * 60 * 60

func SubmitGame(db *db.Db) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		magic := ""
		cookie, err := r.Cookie(submitgame.MagicNumberCookie)
		if err != nil {
			magic = r.URL.Query().Get(submitgame.MagicNumberParam)
			http.SetCookie(w, &http.Cookie{
				Name:   submitgame.MagicNumberCookie,
				Value:  magic,
				MaxAge: maxAge,
			})
		} else {
			magic = cookie.Value
		}

		if magic != magicNumber {
			slog.Warn("An attempt to access the submit form without the magic number was made")
			w.Write([]byte("This page can only be accessed from the QR code in the pub"))
			return
		}

		assignNewIkey := false
		ikeyCookie, err := r.Cookie(submitgame.IKeyCookie)
		if err != nil {
			slog.Info("No ikey for user trying to submit game")
			assignNewIkey = true
		} else {
			ikey, err := strconv.ParseInt(ikeyCookie.Value, 10, 64)
			if err != nil || ikey < 0 {
				slog.Warn("Invalid ikey detected", "err", err, "key", ikeyCookie.Value)
				assignNewIkey = true
			}
		}

		if assignNewIkey {
			ikey, err := model.NextIKey(db)
			if err != nil {
				slog.Warn("Could not assign new ikey", "err", err)
				w.Write([]byte("Could not generate an idempotency - unable to report a game"))
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     submitgame.IKeyCookie,
				Value:    fmt.Sprintf("%d", ikey),
				MaxAge:   maxAge,
				HttpOnly: true,
				Secure:   true,
				Path:     submitgame.BasePath,
			})
		}

		body, err := submitgame.Render()
		if err != nil {
			slog.Error("Cannot render submit game", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		Render(w, body)
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("alive and well at %s", time.Now().UTC())
	Render(w, template.HTML(fmt.Sprintf("<p>%s</p>", msg)))
}

// NewHandler returns a router that handles all site routes.
func NewHandler(db *db.Db) http.Handler {
	mux := http.NewServeMux()
	// {$} matches exactly "/"
	mux.HandleFunc("GET /{$}", Index(db))
	mux.HandleFunc("GET /test", Test)
	mux.HandleFunc("GET /privacy-policy", PrivacyPolicy)

	slog.Info(fmt.Sprintf("To submit a game use http://0.0.0.0:8080/%s?%s=%s",
		submitgame.BasePath,
		submitgame.MagicNumberParam,
		magicNumber))
	mux.HandleFunc(fmt.Sprintf("GET %s", submitgame.BasePath), SubmitGame(db))
	submitgame.Register(mux, db)

	return mux
}
