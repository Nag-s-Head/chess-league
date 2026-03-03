package playerdetails

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed details.html
var f embed.FS
var tpl *template.Template = utils.GetTemplate(f, "details.html")

type PlayerDetails struct {
	Player      model.Player
	TotalGames  int
	Wins        int
	Losses      int
	Draws       int
	WinRate     float64
	LossRate    float64
	DrawRate    float64
	GameHistory []GameWithOutcome
}

type GameWithOutcome struct {
	OpponentName string
	Outcome      string
	Played       time.Time
	EloChange    int
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
		Player:      player,
		TotalGames:  len(games),
		GameHistory: make([]GameWithOutcome, 0),
	}

	for _, g := range games {
		gw := GameWithOutcome{
			Played: g.Played,
		}

		isWhite := g.PlayerWhite == id
		if isWhite {
			gw.OpponentName = g.BlackName
			if g.Score == model.Score_Win {
				gw.Outcome = "Win"
				gw.EloChange = g.EloGiven
				details.Wins++
			} else if g.Score == model.Score_Loss {
				gw.Outcome = "Loss"
				gw.EloChange = g.EloTaken
				details.Losses++
			} else {
				gw.Outcome = "Draw"
				// We don't know who was underdog in draws without more data,
				// so we'll show 0 or handle it gracefully.
				// For simplicity in history, show 0 if not sure.
				gw.EloChange = 0
				details.Draws++
			}
		} else {
			gw.OpponentName = g.WhiteName
			if g.Score == model.Score_Loss {
				gw.Outcome = "Win"
				gw.EloChange = g.EloGiven
				details.Wins++
			} else if g.Score == model.Score_Win {
				gw.Outcome = "Loss"
				gw.EloChange = g.EloTaken
				details.Losses++
			} else {
				gw.Outcome = "Draw"
				gw.EloChange = 0
				details.Draws++
			}
		}
		details.GameHistory = append(details.GameHistory, gw)
	}

	if details.TotalGames > 0 {
		details.WinRate = float64(details.Wins) / float64(details.TotalGames) * 100
		details.LossRate = float64(details.Losses) / float64(details.TotalGames) * 100
		details.DrawRate = float64(details.Draws) / float64(details.TotalGames) * 100
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, details)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return template.HTML(buf.String()), nil
}
