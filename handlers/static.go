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

	Render(w, body)
	utils.WithCacheControl(w)
}

func Rules(w http.ResponseWriter, r *http.Request) {
	body, err := rules.Render()
	if err != nil {
		slog.Error("Cannot render rules", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	Render(w, body)
	utils.WithCacheControl(w)
}
