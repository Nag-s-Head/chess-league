package rules

import (
	"bytes"
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
	githubapi "github.com/Nag-s-Head/chess-league/handlers/utils/github_api"
)

const (
	CurrentRulesVersion = "1"
	RulesVersionCookie  = "rules-version"
)

func HasAgreedToRules(r *http.Request) bool {
	cookie, err := r.Cookie(RulesVersionCookie)
	if err != nil {
		return false
	}
	return cookie.Value == CurrentRulesVersion
}

//go:embed rules.html
var f embed.FS
var policy *template.Template = utils.GetTemplate(f, "rules.html")

type RulesData struct {
	Members         []githubapi.User
	ShowAgreeButton bool
}

func Render(showAgreeButton bool) (template.HTML, error) {
	var members []githubapi.User = []githubapi.User{}
	members, err := githubapi.GerOrganisationMembers(os.Getenv("GITHUB_ORGANISATION"), os.Getenv("GITHUB_API_KEY"))
	if err != nil {
		slog.Warn("Was not able to get organisation members", "err", err)
	}

	data := RulesData{
		Members:         members,
		ShowAgreeButton: showAgreeButton,
	}

	var buf bytes.Buffer
	err = policy.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
