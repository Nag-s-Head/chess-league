package elo

import (
	"math"

	"github.com/Nag-s-Head/chess-league/db/model"
)

const K = 30

func p(a, b int) float64 {
	return 1.0 / (1.0 + float64(math.Pow(10.0, (float64(a-b)/400.0))))
}

type Outcome float64

const (
	Outcome_Win  = 1.0
	Outcome_Loss = 0.0
	Outcome_Draw = 0.5
)

// Updates player A, and B's ELO based on the outcome, see Outcome_XXX.
// The outcome describes player a, so if it is Outcome_Win then
func CalculateElo(a, b *model.Player, outcome Outcome) (int, int) {
	outcomeF := float64(outcome)
	pb := p(a.Elo, b.Elo)
	pa := p(b.Elo, a.Elo)

	deltaA := int(math.Round(K * (outcomeF - pa)))
	a.Elo += deltaA

	deltaB := int(math.Round(K * ((1 - outcomeF) - pb)))
	b.Elo += deltaB

	return deltaA, deltaB
}
