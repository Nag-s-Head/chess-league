package players

import (
	"embed"
	"html/template"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed players.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "players.html")
