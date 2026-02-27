package model_test 

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/stretchr/testify/require"
)

func TestEloWin(t *testing.T) {
	a := model.NewPlayer("Dave")
	b := model.NewPlayer("Bob")

	deltaA, deltaB := model.CalculateElo(&a, &b, model.Outcome_Win)
	require.Equal(t, model.StartingElo+deltaA, a.Elo)
	require.Equal(t, model.StartingElo+deltaB, b.Elo)

	require.Equal(t, deltaA, 15)
	require.Equal(t, deltaB, -deltaA)
}

func TestEloLoss(t *testing.T) {
	a := model.NewPlayer("Dave")
	b := model.NewPlayer("Bob")

	deltaA, deltaB := model.CalculateElo(&a, &b, model.Outcome_Loss)
	require.Equal(t, model.StartingElo+deltaA, a.Elo)
	require.Equal(t, model.StartingElo+deltaB, b.Elo)

	require.Equal(t, deltaA, -15)
	require.Equal(t, deltaB, -deltaA)
}

func TestEloDraw(t *testing.T) {
	a := model.NewPlayer("Dave")
	b := model.NewPlayer("Bob")

	deltaA, deltaB := model.CalculateElo(&a, &b, model.Outcome_Draw)
	require.Equal(t, model.StartingElo+deltaA, a.Elo)
	require.Equal(t, model.StartingElo+deltaB, b.Elo)

	require.Equal(t, deltaA, 0)
	require.Equal(t, deltaB, -deltaA)
}
