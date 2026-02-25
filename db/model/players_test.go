package model_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewPlayer(t *testing.T) {
	name := uuid.New().String()
	player := model.NewPlayer(name)

	require.NotEmpty(t, player.Id)
	require.Equal(t, name, player.Name)
	require.Equal(t, name, player.NameNormalised)
	require.NotEmpty(t, player.JoinTime)
	require.Equal(t, model.StartingElo, player.Elo)
}

func TestInsertPlayer(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	player := model.NewPlayer(name)

	require.NoError(t, model.InsertPlayer(db, player))
}

func TestInsertLotsOfPlayers(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	for range 100 {
		name := uuid.New().String()
		player := model.NewPlayer(name)
		require.NoError(t, model.InsertPlayer(db, player))
	}
}

func TestGetPlayer(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	player := model.NewPlayer(name)

	require.NoError(t, model.InsertPlayer(db, player))

	player2, err := model.GetPlayer(db, player.Id)

	require.InDelta(t, player.JoinTime.Unix(), player2.JoinTime.Unix(), 1)
	// Little hack for the time zone
	player2.JoinTime = player.JoinTime

	require.NoError(t, err)
	require.Equal(t, player, player2)
}
