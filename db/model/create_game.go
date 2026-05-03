package model

import (
	"errors"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

func CreateGame(tx *sqlx.Tx, submitter, opponent *Player, submitterIsWhite bool, ikey int64, score Score, r *http.Request) (Game, int, int, error) {
	if submitter.Id == opponent.Id {
		return Game{}, 0, 0, errors.New("Both players are the same")
	}

	game := Game{
		Score:           score,
		Submitter:       submitter.Id,
		Played:          time.Now(),
		Deleted:         false,
		SubmitIp:        GetRemoteAddr(r),
		SubmitUserAgent: r.UserAgent(),
		IKey:            ikey,
	}

	var pWhite, pBlack *Player
	if submitterIsWhite {
		pWhite = submitter
		pBlack = opponent
	} else {
		pWhite = opponent
		pBlack = submitter
	}

	var outcome Outcome
	switch score {
	case Score_Win:
		outcome = Outcome_Win
	case Score_Loss:
		outcome = Outcome_Loss
	case Score_Draw:
		outcome = Outcome_Draw
	}

	game.Liglicko2WhiteOldRating = pWhite.Liglicko2Rating
	game.Liglicko2WhiteOldVolatility = pWhite.Liglicko2Volatility
	game.Liglicko2WhiteOldDeviation = pWhite.Liglicko2Deviation
	game.Liglicko2WhiteOldAt = pWhite.Liglicko2At

	game.Liglicko2BlackOldRating = pBlack.Liglicko2Rating
	game.Liglicko2BlackOldVolatility = pBlack.Liglicko2Volatility
	game.Liglicko2BlackOldDeviation = pBlack.Liglicko2Deviation
	game.Liglicko2BlackOldAt = pBlack.Liglicko2At

	eloWhite, eloBlack := CalculateElo(pWhite, pBlack, outcome)
	liglicko2White, liglicko2Black, err := CalculateLiglicko2(pWhite, pBlack, outcome, game.Played)
	if err != nil {
		return Game{}, 0, 0, errors.Join(errors.New("Could not calculate liglicko2"), err)
	}
	game.PlayerWhite = pWhite.Id
	game.PlayerBlack = pBlack.Id

	if eloWhite > eloBlack {
		game.DEPRECATEDEloGiven = eloWhite
		game.DEPRECATEDEloTaken = eloBlack
	} else {
		game.DEPRECATEDEloGiven = eloBlack
		game.DEPRECATEDEloTaken = eloWhite
	}

	game.Liglicko2White = liglicko2White
	game.Liglicko2Black = liglicko2Black

	_, err = tx.NamedExec(`
INSERT INTO games (
	player_white, player_black, score, submitter, played, deleted, elo_given, elo_taken, liglicko2_white, liglicko2_black, submit_ip, submit_user_agent, ikey,
	liglicko2_white_old_rating, liglicko2_white_old_volatility, liglicko2_white_old_deviation, liglicko2_white_old_at,
	liglicko2_black_old_rating, liglicko2_black_old_volatility, liglicko2_black_old_deviation, liglicko2_black_old_at
)
VALUES (
	:player_white, :player_black, :score, :submitter, :played, :deleted, :elo_given, :elo_taken, :liglicko2_white, :liglicko2_black, :submit_ip, :submit_user_agent, :ikey,
	:liglicko2_white_old_rating, :liglicko2_white_old_volatility, :liglicko2_white_old_deviation, :liglicko2_white_old_at,
	:liglicko2_black_old_rating, :liglicko2_black_old_volatility, :liglicko2_black_old_deviation, :liglicko2_black_old_at
);
  	`, game)

	if err != nil {
		return Game{}, 0, 0, errors.Join(errors.New("Cannot insert game"), err)
	}

	_, err = tx.NamedExec(`UPDATE players 
SET elo=:elo, liglicko2_rating=:liglicko2_rating, liglicko2_deviation=:liglicko2_deviation, liglicko2_volatility=:liglicko2_volatility, liglicko2_at=:liglicko2_at
WHERE id=:id`, pWhite)
	if err != nil {
		return Game{}, 0, 0, errors.Join(errors.New("Cannot set elo of white player"), err)
	}

	_, err = tx.NamedExec(`UPDATE players 
SET elo=:elo, liglicko2_rating=:liglicko2_rating, liglicko2_deviation=:liglicko2_deviation, liglicko2_volatility=:liglicko2_volatility, liglicko2_at=:liglicko2_at
WHERE id=:id`, pBlack)
	if err != nil {
		return Game{}, 0, 0, errors.Join(errors.New("Cannot set elo of black player"), err)
	}

	return game, int(liglicko2White), int(liglicko2Black), nil
}
