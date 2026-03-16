package rules

import (
	"bytes"
	"embed"
	"html/template"
	"log/slog"
	"os"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
	githubapi "github.com/Nag-s-Head/chess-league/handlers/utils/github_api"
)

//go:embed rules.html
var f embed.FS
var policy *template.Template = utils.GetTemplate(f, "rules.html")

func Render() (template.HTML, error) {
	var members []githubapi.User = []githubapi.User{}
	members, err := githubapi.GerOrganisationMembers(os.Getenv("GITHUB_ORGANISATION"))
	if err != nil {
		slog.Warn("Was not able to get organisation members", "err", err)
	}

	var buf bytes.Buffer
	err = policy.Execute(&buf, members)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
