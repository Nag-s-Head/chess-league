package auth

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Nag-s-Head/chess-league/db"
)

const AuthCookie = "admin-authentication"

var isTestMode = os.Getenv("TEST_MODE") == "true"

func WithAuthentication(db *db.Db, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AuthCookie)
		if err != nil {
			slog.Info("A user has tried to access the admin portal without being logged in, redirecting to authentication page")

			url := "https://github.com/CHANGE ME"
			if isTestMode {
				url = "/admin/test-mode"
			}

			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		slog.Error("TODO: handle the authentication cookie", "cookie", cookie)

		// TODO: check the session id stored in cookie
		next(w, r)
	}
}
