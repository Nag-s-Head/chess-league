package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	githubapi "github.com/Nag-s-Head/chess-league/handlers/utils/github_api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var oauthConfig *oauth2.Config

func init() {
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/admin/auth/callback", baseURL),
		Scopes:       []string{"read:user", "read:org"},
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if isTestMode {
		http.Redirect(w, r, "/admin/test-mode", http.StatusTemporaryRedirect)
		return
	}

	url := oauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func Callback(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			slog.Error("No code in github callback")
			http.Error(w, "No code in callback", http.StatusBadRequest)
			return
		}

		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			slog.Error("Could not exchange code for token", "err", err)
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		ghUser, err := githubapi.GetAuthenticatedUser(token.AccessToken)
		if err != nil {
			slog.Error("Could not get authenticated user from Github", "err", err)
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}

		orgName := os.Getenv("GITHUB_ORGANISATION")
		apiKey := os.Getenv("GITHUB_API_KEY")

		isMember, err := githubapi.IsMemberOfOrg(orgName, ghUser.Login, apiKey)
		if err != nil {
			slog.Error("Could not check org membership", "err", err, "org", orgName, "user", ghUser.Login)
			http.Error(w, "Failed to verify organisation membership", http.StatusInternalServerError)
			return
		}

		if !isMember {
			slog.Warn("User is not a member of the required organisation", "user", ghUser.Login, "org", orgName)
			http.Error(w, "You are not an authorized admin", http.StatusForbidden)
			return
		}

		// Login successful, create or update admin user in DB
		adminUser, err := model.AdminLogin(db, ghUser.Name, ghUser.Login, r.RemoteAddr, r.UserAgent())
		if err != nil {
			slog.Error("Could not login admin user in database", "err", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, CreateAuthCookie(adminUser.SessionKey))
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
	}
}
