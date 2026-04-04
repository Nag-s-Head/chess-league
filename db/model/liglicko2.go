package model

import (
	"time"

	liglicko2core "github.com/Nag-s-Head/chess-league/db/liglicko2"
)

type Liglicko2Rating = liglicko2core.Rating

func liglicko2InstantFromTime(t time.Time) float64 {
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
	now := liglicko2InstantFromTime(playedAt)
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
