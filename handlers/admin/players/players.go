package players

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed players.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "players.html")

func Render(db *db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		users, err := model.GetPlayersByEloWithGameCount(db)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		err = indexTpl.Execute(&buf, users)
		if err != nil {
			return "", err
		}
		return template.HTML(buf.String()), nil
	}
}
