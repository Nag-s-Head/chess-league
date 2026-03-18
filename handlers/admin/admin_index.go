package admin

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed admin_index.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "admin_index.html")

func AdminIndex(w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	var buf bytes.Buffer
	err := indexTpl.Execute(&buf, nil)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
