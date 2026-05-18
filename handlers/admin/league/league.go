package league

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed league.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "league.html")

func Render(db *db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		players, err := model.GetUiFriendlyLeaguePlayers(db)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		err = indexTpl.Execute(&buf, players)
		if err != nil {
			return "", err
		}

		return template.HTML(buf.String()), nil
	}
}
