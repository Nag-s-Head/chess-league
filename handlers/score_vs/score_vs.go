package scorevs

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed *.html
var templates embed.FS

var scoreVsTemplate = utils.GetTemplate(templates, "score_vs.html")

type ScoreVsData struct {
	Player   model.Player
	Opponent model.Player
	Wins     int
	Draws    int
	Losses   int
	Games    []model.GameWithOutcome
}

func RenderScoreVs(dbCon db.Db, playerID, opponentID uuid.UUID) (template.HTML, error) {
	if playerID == opponentID {
		return "", errors.New("a player cannot play against themselves")
	}

	player, err := model.GetPlayer(dbCon, playerID)
	if err != nil {
		return "", fmt.Errorf("cannot get player: %w", err)
	}
	opponent, err := model.GetPlayer(dbCon, opponentID)
	if err != nil {
		return "", fmt.Errorf("cannot get opponent: %w", err)
	}

	playerGames, err := model.GetGamesByPlayer(dbCon, playerID)
	if err != nil {
		return "", fmt.Errorf("cannot get player games: %w", err)
	}

	data := ScoreVsData{Player: player, Opponent: opponent, Games: make([]model.GameWithOutcome, 0)}
	for _, game := range playerGames {
		if game.Deleted || !((game.PlayerWhite == playerID && game.PlayerBlack == opponentID) || (game.PlayerBlack == playerID && game.PlayerWhite == opponentID)) {
			continue
		}

		outcome := game.MapGameToGameWithOutcome(playerID)
		data.Games = append(data.Games, outcome)
		switch outcome.Outcome {
		case "Win":
			data.Wins++
		case "Draw":
			data.Draws++
		case "Loss":
			data.Losses++
		}
	}

	return executeTemplate(scoreVsTemplate, data)
}

func executeTemplate(tmpl *template.Template, data any) (template.HTML, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("cannot execute score-vs template: %w", err)
	}
	return template.HTML(buf.String()), nil
}
