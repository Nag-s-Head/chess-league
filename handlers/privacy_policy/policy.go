package privacypolicy

import (
	"bytes"
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed policy.html
var f embed.FS
var policy *template.Template = utils.GetTemplate(f, "policy.html")

func Render() (template.HTML, error) {
	var buf bytes.Buffer
	err := policy.Execute(&buf, nil)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func Policy(w http.ResponseWriter, r *http.Request) {
	err := policy.Execute(w, nil)
	if err != nil {
		slog.Error("Cannot execute policy template", "err", err)
	}
}
