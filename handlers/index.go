package handlers

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed index.html
var f embed.FS
var tpl *template.Template = utils.GetTemplate(f, "index.html")

func Index(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, nil)
	if err != nil {
		slog.Error("Cannot execute index template", "err", err)
	}
}
