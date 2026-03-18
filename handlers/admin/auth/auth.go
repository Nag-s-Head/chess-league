package auth

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
)

const AuthCookie = "admin-authentication"

var isTestMode = os.Getenv("TEST_MODE") == "true"

// If the key is empty string it will remove the cookie
func CreateAuthCookie(sessionKey string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     AuthCookie,
		Secure:   true,
		HttpOnly: true,
		MaxAge:   3600,
		Value:    sessionKey,
	}

	if sessionKey == "" {
		cookie.MaxAge = 0
	}

	return cookie
}

func loginUrl() string {
	if isTestMode {
		return "/admin/test-mode"
	}

	return "https://github.com/CHANGE ME"
}

func WithAuthentication(db *db.Db, next func(w http.ResponseWriter, r *http.Request, user *model.AdminUser)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AuthCookie)
		if err != nil {
			slog.Info("A user has tried to access the admin portal without being logged in, redirecting to authentication page")

			http.Redirect(w, r, loginUrl(), http.StatusTemporaryRedirect)
			return
		}

		user, err := model.AdminGetFromSessionKey(db, cookie.Value)
		if err != nil {
			slog.Warn("User with invalid authentication tried to access the page", "url", r.URL, "err", err)

			http.SetCookie(w, CreateAuthCookie(""))
			http.Redirect(w, r, loginUrl(), http.StatusTemporaryRedirect)
			return
		}
		next(w, r, user)
	}
}
