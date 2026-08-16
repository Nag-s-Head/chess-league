package playerdetails

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	elo_charts "github.com/Nag-s-Head/chess-league/img/elo_charts"
	"github.com/google/uuid"
)

//go:embed details.html
var f embed.FS
var tpl *template.Template = utils.GetTemplate(f, "details.html")

type PlayerDetails struct {
	Player  model.Player
	Details model.GamesUiFriendly
}

func Render(dbCon db.Db, id uuid.UUID) (template.HTML, error) {
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

func ServeChart(db db.Db) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			slog.Warn("Invalid player ID", "id", idStr)
			http.Error(w, "Invalid player ID", http.StatusBadRequest)
			return
		}

		chart, err := model.PlayerEloChart(db, id)
		if err != nil {
			slog.Warn("Invalid player ID", "id", idStr)
			http.Error(w, "Invalid player ID", http.StatusBadRequest)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=300, s-max-age=300")
		w.Header().Add("Content-Type", elo_charts.ContentType)
		w.Write(chart)
	}
}
