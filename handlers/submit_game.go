package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	submitgame "github.com/Nag-s-Head/chess-league/handlers/submit_game"
)

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
