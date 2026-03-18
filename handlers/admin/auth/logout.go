package auth

import (
	"log/slog"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
)

func Logout(db *db.Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		cookie, err := r.Cookie(AuthCookie)
		if err != nil {
			return
		}

		user, err := model.AdminGetFromSessionKey(db, cookie.Value)
		if err != nil {
			slog.Error("Could not check the user session key", "err", err)
			return
		}

		err = model.AdminLogout(db, user.Id)
		if err != nil {
			slog.Error("Could not check the user session key", "err", err)
			return
		}

		slog.Info("Logged a user out", "id", user.Id, "name", user.Name, "oauth id", user.OauthId)
	}
}
