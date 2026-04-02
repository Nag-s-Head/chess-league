package player_details

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed player_details.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "player_details.html")

type Model struct {
	Player  model.Player
	Details model.GamesUiFriendly
}

func Render(db *db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return "", err
		}

		player, err := model.GetPlayer(db, id)
		if err != nil {
			return "", err
		}

		games, err := model.GetGamesByPlayer(db, id)
		if err != nil {
			return "", err
		}

		details := model.MapGamesToUserFriendly(id, games)

		var buf bytes.Buffer
		err = indexTpl.Execute(&buf, Model{
			Player:  player,
			Details: details,
		})
		if err != nil {
			return "", err
		}
		return template.HTML(buf.String()), nil
	}
}
