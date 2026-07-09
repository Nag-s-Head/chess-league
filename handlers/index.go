package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Nag-s-Head/chess-league/app/theme"
	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
)

type IndexData struct {
	Players                      []model.Player
	TotalGames                   int
	TotalPlayers                 int
	MinimumStableRatingDeviation float64
	Theme                        theme.Theme
}

func Index(db db.Db, theme theme.Theme) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		players, err := model.GetPlayersByElo(db, false)
		if err != nil {
			slog.Warn("Could not get leaderboard", "err", err)
		}

		gameCount, err := model.GetTotalGameCount(db)
		if err != nil {
			slog.Warn("Could not get game count", "err", err)
		}

		playerCount, err := model.GetTotalPlayerCount(db)
		if err != nil {
			slog.Warn("Could not get player count", "err", err)
		}

		data := IndexData{
			Players:                      players,
			TotalGames:                   gameCount,
			TotalPlayers:                 playerCount,
			MinimumStableRatingDeviation: model.MinimumStableRatingDeviation,
		}

		var buf bytes.Buffer
		err = indexTmpl.Execute(&buf, data)
		if err != nil {
			slog.Error("Cannot execute index template", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		WithLayout(theme)(w, template.HTML(buf.String()))
	}
}
