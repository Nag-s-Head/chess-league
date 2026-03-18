package admin

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/handlers/admin/auth"
	testmode "github.com/Nag-s-Head/chess-league/handlers/admin/test_mode"
)

const (
	BasePath = "/admin"
)

type PageRenderer func(w http.ResponseWriter, r *http.Request) (template.HTML, error)
type LayoutRenderer func(w http.ResponseWriter, body template.HTML)

var isTestMode = os.Getenv("TEST_MODE") == "true"

func WithLayout(Render PageRenderer, LayoutRender LayoutRenderer) func(w http.ResponseWriter, r *http.Request) {
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

func Register(mux *http.ServeMux, db *db.Db, LayoutRender func(w http.ResponseWriter, body template.HTML)) {
	if isTestMode {
		slog.Warn("Test mod is enabled, if this is a production environment then you should turn it off!")
		mux.HandleFunc(fmt.Sprintf("GET %s/test-mode", BasePath), WithLayout(testmode.Login, LayoutRender))
		mux.HandleFunc(fmt.Sprintf("POST %s/test-mode", BasePath), testmode.LoginPost)
	}

	mux.HandleFunc(fmt.Sprintf("GET %s", BasePath), auth.WithAuthentication(db, WithLayout(AdminIndex, LayoutRender)))
}
