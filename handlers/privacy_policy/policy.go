package privacypolicy

import (
	"bytes"
	"embed"
	"html/template"

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
