package handlers

import (
	"log/slog"
	"net/http"

	privacypolicy "github.com/Nag-s-Head/chess-league/handlers/privacy_policy"
	"github.com/Nag-s-Head/chess-league/handlers/rules"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

func PrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	body, err := privacypolicy.Render()
	if err != nil {
		slog.Error("Cannot render privacy policy", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	utils.WithCacheControl(w)
	WithLayout(w, body)
}

func Rules(w http.ResponseWriter, r *http.Request) {
	showAgreeButton := r.URL.Query().Get("agree") == "true"
	body, err := rules.Render(showAgreeButton)
	if err != nil {
		slog.Error("Cannot render rules", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	utils.WithCacheControl(w)
	WithLayout(w, body)
}

func RulesAgree(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     rules.RulesVersionCookie,
		Value:    rules.CurrentRulesVersion,
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	http.Redirect(w, r, "/submit-game", http.StatusFound)
}
