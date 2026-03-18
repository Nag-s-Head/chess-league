package admin

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
)

const (
	BasePath   = "/admin"
	AuthCookie = "admin-authentication"
)

type PageRenderer func(w http.ResponseWriter, r *http.Request) (template.HTML, error)
type LayoutRenderer func(w http.ResponseWriter, body template.HTML)

func WithLayout(Render PageRenderer, LayoutRender LayoutRenderer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := Render(w, r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Could not render page: %w", err)))
			slog.Error("Could not render admin portal page", "err", err, "url", r.URL)
			return
		}

		LayoutRender(w, tpl)
	}
}

func Register(mux *http.ServeMux, db *db.Db, LayoutRender func(w http.ResponseWriter, body template.HTML)) {
	mux.HandleFunc(fmt.Sprintf("GET %s", BasePath), WithLayout(AdminIndex, LayoutRender))
}
