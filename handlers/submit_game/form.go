package submitgame

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed form.html submit.html
var f embed.FS
var tpl *template.Template = utils.GetTemplate(f, "form.html")
var resultsTpl *template.Template = utils.GetTemplate(f, "submit.html")

func Render() (template.HTML, error) {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, nil)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
