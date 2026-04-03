package model_test

import (
	"testing"
	"time"

	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/stretchr/testify/require"
)

func TestLiglicko2Win(t *testing.T) {
	a := model.NewPlayer("Dave")
	b := model.NewPlayer("Bob")

	deltaA, deltaB, err := model.CalculateLiglicko2(&a, &b, model.Outcome_Win, time.Now())
	require.NoError(t, err)
	require.Greater(t, deltaA, 0.0)
	require.Less(t, deltaB, 0.0)
	require.Greater(t, a.Liglicko2Rating, b.Liglicko2Rating)
}

func TestLiglicko2DrawEqualPlayers(t *testing.T) {
	a := model.NewPlayer("Dave")
	b := model.NewPlayer("Bob")

	deltaA, deltaB, err := model.CalculateLiglicko2(&a, &b, model.Outcome_Draw, time.Now())
	require.NoError(t, err)
	require.InDelta(t, 0.0, deltaA, 1e-9)
	require.InDelta(t, 0.0, deltaB, 1e-9)
}
