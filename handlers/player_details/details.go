package playerdetails

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed details.html
var f embed.FS
var tpl *template.Template = utils.GetTemplate(f, "details.html")

type PlayerDetails struct {
	Player  model.Player
	Details model.GamesUiFriendly
}

func Render(dbCon *db.Db, id uuid.UUID) (template.HTML, error) {
	player, err := model.GetPlayer(dbCon, id)
	if err != nil {
		return "", err
	}

	games, err := model.GetGamesByPlayer(dbCon, id)
	if err != nil {
		return "", err
	}

	details := PlayerDetails{
		Player:  player,
		Details: model.MapGamesToUserFriendly(id, games),
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, details)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return template.HTML(buf.String()), nil
}
