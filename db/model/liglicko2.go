package model

import (
	"errors"
	"fmt"
	"time"

	liglicko2core "github.com/Nag-s-Head/chess-league/db/liglicko2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Liglicko2Rating = liglicko2core.Rating

func Liglicko2InstantFromTime(t time.Time) float64 {
	return liglicko2core.InstantFromTime(t)
}

func liglicko2ScoreFromOutcome(outcome Outcome) float64 {
	return liglicko2core.Clamp(float64(outcome), 0.0, 1.0)
}

func playerLiglicko2Rating(p *Player) Liglicko2Rating {
	return Liglicko2Rating{
		Rating:     p.Liglicko2Rating,
		Deviation:  p.Liglicko2Deviation,
		Volatility: p.Liglicko2Volatility,
		At:         p.Liglicko2At,
	}
}

func setPlayerLiglicko2Rating(p *Player, r Liglicko2Rating) {
	p.Liglicko2Rating = r.Rating
	p.Liglicko2Deviation = r.Deviation
	p.Liglicko2Volatility = r.Volatility
	p.Liglicko2At = r.At
}

func CalculateLiglicko2(a, b *Player, outcome Outcome, playedAt time.Time) (float64, float64, error) {
	now := Liglicko2InstantFromTime(playedAt)
	first := playerLiglicko2Rating(a)
	second := playerLiglicko2Rating(b)
	score := liglicko2ScoreFromOutcome(outcome)

	nextFirst, err := liglicko2core.UpdateSingle(first, second, score, now, liglicko2core.FirstAdvantage)
	if err != nil {
		return 0, 0, err
	}

	nextSecond, err := liglicko2core.UpdateSingle(second, first, 1.0-score, now, -liglicko2core.FirstAdvantage)
	if err != nil {
		return 0, 0, err
	}

	deltaA := nextFirst.Rating - first.Rating
	deltaB := nextSecond.Rating - second.Rating

	setPlayerLiglicko2Rating(a, nextFirst)
	setPlayerLiglicko2Rating(b, nextSecond)

	return deltaA, deltaB, nil
}

func ReplayFrom(txx *sqlx.Tx, ikey int64) ([]Game, []*Player, error) {
	getOrAddPlayer := func(txx *sqlx.Tx, players map[uuid.UUID]*Player, id uuid.UUID, game *Game) (*Player, error) {
		player, found := players[id]
		if !found {
			p, err := getPlayerById(txx, id)
			if err != nil {
				return nil, err
			}

			if game.PlayerWhite == id {
				p.Liglicko2Rating = game.Liglicko2WhiteOldRating
				p.Liglicko2Deviation = game.Liglicko2WhiteOldDeviation
				p.Liglicko2Volatility = game.Liglicko2WhiteOldVolatility
				p.Liglicko2At = game.Liglicko2WhiteOldAt
			} else {
				p.Liglicko2Rating = game.Liglicko2BlackOldRating
				p.Liglicko2Deviation = game.Liglicko2BlackOldDeviation
				p.Liglicko2Volatility = game.Liglicko2BlackOldVolatility
				p.Liglicko2At = game.Liglicko2BlackOldAt
			}

			player = &p
			players[id] = player
		}

		return player, nil
	}

	var seedGame Game
	affectedPlayers := make(map[uuid.UUID]*Player)
	err := txx.Get(&seedGame, "SELECT * FROM games WHERE ikey = $1;", ikey)
	if err != nil {
		return nil, nil, errors.Join(errors.New("Cannot get seed game to start replays from"), err)
	}

	// Force load the old state of the players even if the game is deleted
	_, err = getOrAddPlayer(txx, affectedPlayers, seedGame.PlayerWhite, &seedGame)
	_, err = getOrAddPlayer(txx, affectedPlayers, seedGame.PlayerBlack, &seedGame)

	var affectedGames []Game
	err = txx.Select(&affectedGames, "SELECT * FROM games WHERE played >= $1 AND deleted = FALSE ORDER BY played ASC;", seedGame.Played)
	if err != nil {
		return nil, nil, errors.Join(errors.New("Cannot select affected players"), err)
	}

	for i, game := range affectedGames {
		white, err := getOrAddPlayer(txx, affectedPlayers, game.PlayerWhite, &game)
		if err != nil {
			return nil, nil, errors.Join(errors.New("Cannot get white player"), err)
		}

		black, err := getOrAddPlayer(txx, affectedPlayers, game.PlayerBlack, &game)
		if err != nil {
			return nil, nil, errors.Join(errors.New("Cannot get black player"), err)
		}

		whiteRating := liglicko2core.Rating{
			Rating:     white.Liglicko2Rating,
			Deviation:  white.Liglicko2Deviation,
			Volatility: white.Liglicko2Volatility,
			At:         white.Liglicko2At,
		}

		blackRating := liglicko2core.Rating{
			Rating:     black.Liglicko2Rating,
			Deviation:  black.Liglicko2Deviation,
			Volatility: black.Liglicko2Volatility,
			At:         black.Liglicko2At,
		}

		now := liglicko2core.InstantFromTime(game.Played)
		score := liglicko2ScoreFromOutcome(game.Score.Outcome())

		nextWhite, err := liglicko2core.UpdateSingle(whiteRating, blackRating, score, now, liglicko2core.FirstAdvantage)
		if err != nil {
			return nil, nil, errors.Join(fmt.Errorf("Cannot replay white rating for game %d", game.IKey), err)
		}

		nextBlack, err := liglicko2core.UpdateSingle(blackRating, whiteRating, 1.0-score, now, -liglicko2core.FirstAdvantage)
		if err != nil {
			return nil, nil, errors.Join(fmt.Errorf("Cannot replay black rating for game %d", game.IKey), err)
		}

		affectedGames[i].Liglicko2White = nextWhite.Rating - whiteRating.Rating
		affectedGames[i].Liglicko2Black = nextBlack.Rating - blackRating.Rating
		affectedGames[i].Liglicko2WhiteOldRating = whiteRating.Rating
		affectedGames[i].Liglicko2WhiteOldVolatility = whiteRating.Volatility
		affectedGames[i].Liglicko2WhiteOldDeviation = whiteRating.Deviation
		affectedGames[i].Liglicko2WhiteOldAt = whiteRating.At
		affectedGames[i].Liglicko2BlackOldRating = blackRating.Rating
		affectedGames[i].Liglicko2BlackOldVolatility = blackRating.Volatility
		affectedGames[i].Liglicko2BlackOldDeviation = blackRating.Deviation
		affectedGames[i].Liglicko2BlackOldAt = blackRating.At

		white.ApplyRating(nextWhite)
		black.ApplyRating(nextBlack)

		_, err = txx.NamedExec(`
UPDATE games SET 
    liglicko2_white = :liglicko2_white,
    liglicko2_black = :liglicko2_black,
    liglicko2_white_old_rating = :liglicko2_white_old_rating,
    liglicko2_white_old_volatility = :liglicko2_white_old_volatility,
    liglicko2_white_old_deviation = :liglicko2_white_old_deviation,
    liglicko2_white_old_at = :liglicko2_white_old_at,
    liglicko2_black_old_rating = :liglicko2_black_old_rating,
    liglicko2_black_old_volatility = :liglicko2_black_old_volatility,
    liglicko2_black_old_deviation = :liglicko2_black_old_deviation,
    liglicko2_black_old_at = :liglicko2_black_old_at
WHERE ikey = :ikey`, affectedGames[i])
		if err != nil {
			return nil, nil, errors.Join(fmt.Errorf("Cannot update game %d during replay", game.IKey), err)
		}
	}

	players := make([]*Player, 0)
	for _, player := range affectedPlayers {
		players = append(players, player)
		_, err := txx.NamedExec(`
UPDATE players SET 
    liglicko2_rating = :liglicko2_rating, 
    liglicko2_deviation = :liglicko2_deviation, 
    liglicko2_volatility = :liglicko2_volatility, 
    liglicko2_at = :liglicko2_at 
WHERE id = :id`, player)
		if err != nil {
			return nil, nil, errors.Join(fmt.Errorf("Cannot update player %s during replay", player.Id), err)
		}
	}

	return affectedGames, players, nil
}
