package admin

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
	githubapi "github.com/Nag-s-Head/chess-league/handlers/utils/github_api"
)

//go:embed admin_index.html
var f embed.FS
var policy *template.Template = utils.GetTemplate(f, "admin_index.html")

type RulesData struct {
	Members         []githubapi.User
	ShowAgreeButton bool
}

func AdminIndex(w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	var buf bytes.Buffer
	err := policy.Execute(&buf, nil)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
