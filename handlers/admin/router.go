package admin

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	adminusers "github.com/Nag-s-Head/chess-league/handlers/admin/admin_users"
	adminuserdetails "github.com/Nag-s-Head/chess-league/handlers/admin/admin_users/admin_user_details"
	"github.com/Nag-s-Head/chess-league/handlers/admin/auth"
	testmode "github.com/Nag-s-Head/chess-league/handlers/admin/test_mode"
)

const (
	BasePath = "/admin"
)

type PageRenderer func(w http.ResponseWriter, r *http.Request) (template.HTML, error)
type PageRendererWithAuth func(w http.ResponseWriter, r *http.Request, user *model.AdminUser) (template.HTML, error)
type LayoutRenderer func(w http.ResponseWriter, body template.HTML)

func WithLayout(Render PageRenderer, LayoutRender LayoutRenderer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := Render(w, r)
		if err != nil {
			w.Write(fmt.Appendf(nil, "Could not render page: %s", err))
			slog.Error("Could not render admin portal page", "err", err, "url", r.URL)
			return
		}

		LayoutRender(w, tpl)
	}
}

func WithLayoutAndAuthentication(db *db.Db, Render PageRendererWithAuth, LayoutRender LayoutRenderer) func(http.ResponseWriter, *http.Request) {
	return auth.WithAuthentication(db, func(w http.ResponseWriter, r *http.Request, user *model.AdminUser) {
		tpl, err := Render(w, r, user)
		if err != nil {
			w.Write(fmt.Appendf(nil, "Could not render page: %s", err))
			slog.Error("Could not render admin portal page", "err", err, "url", r.URL)
			return
		}

		LayoutRender(w, tpl)
	})
}

var isTestMode = os.Getenv("TEST_MODE") == "true"

func Register(mux *http.ServeMux, db *db.Db, LayoutRender func(w http.ResponseWriter, body template.HTML)) {
	if isTestMode {
		slog.Warn("Test mod is enabled, if this is a production environment then you should turn it off!")
		mux.HandleFunc(fmt.Sprintf("GET %s/test-mode", BasePath), WithLayout(testmode.Login, LayoutRender))
		mux.HandleFunc(fmt.Sprintf("POST %s/test-mode", BasePath), testmode.LoginPost(db))
	}

	mux.HandleFunc(fmt.Sprintf("GET %s", BasePath), WithLayoutAndAuthentication(db, AdminIndex, LayoutRender))

	// Auth
	mux.HandleFunc(fmt.Sprintf("GET %s/login", BasePath), auth.Login)
	mux.HandleFunc(fmt.Sprintf("GET %s/auth/callback", BasePath), auth.Callback(db))
	mux.HandleFunc(fmt.Sprintf("GET %s/logout", BasePath), auth.Logout(db))

	// Pages
	mux.HandleFunc(fmt.Sprintf("GET %s/admins", BasePath), WithLayoutAndAuthentication(db, adminusers.Render(db), LayoutRender))
	mux.HandleFunc(fmt.Sprintf("GET %s/admins/{id}", BasePath), WithLayoutAndAuthentication(db, adminuserdetails.Render(db), LayoutRender))
}
