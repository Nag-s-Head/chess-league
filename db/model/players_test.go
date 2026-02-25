package model_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewPlayer(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	for range 100 {
		name := uuid.New().String()
		player := model.NewPlayer(name)
		require.NoError(t, model.InsertPlayer(db, player))
	}
}

func TestGetPlayer(t *testing.T) {
	t.Parallel()
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

func TestGetPlayers(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	player := model.NewPlayer(name)

	require.NoError(t, model.InsertPlayer(db, player))

	players, err := model.GetPlayers(db)
	require.NoError(t, err)
	require.Greater(t, len(players), 1)
}

func TestNormalise(t *testing.T) {
	require.Equal(t, normalisation.Normalise("DANNY PIPER"), normalisation.Normalise("Dánny piper"))
}

func TestSearchPlayerByNameSimpleCase(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	names := []string{"Danny Piper", "Tony Blair", "Gordon Brown", "David Cameron", "Theresa May", "Boris Johnson"}
	for _, name := range names {
		_, err := db.GetSqlxDb().Exec("DELETE FROM players WHERE name_normalised = $1", normalisation.Normalise(name))
		require.NoError(t, err)
	}

	for _, name := range names {
		require.NoError(t, model.InsertPlayer(db, model.NewPlayer(name)))
	}

	players, err := model.SearchPlayerByName(db, "DANNY")
	require.NoError(t, err)
	require.Len(t, players, 1)
	require.Equal(t, players[0].Name, "Danny Piper")
}

func TestSearchPlayerByNameManyResults(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	names := []string{"Liz Truss", "Rishi Sunak", "Kier Starmer", "Joe Bloggs", "Joe Smith", "Joe Bell"}
	for _, name := range names {
		_, err := db.GetSqlxDb().Exec("DELETE FROM players WHERE name_normalised = $1", normalisation.Normalise(name))
		require.NoError(t, err)
	}

	for _, name := range names {
		require.NoError(t, model.InsertPlayer(db, model.NewPlayer(name)))
	}

	t.Run("1 JOE expected", func(t *testing.T) {
		players, err := model.SearchPlayerByName(db, "JOE BLO")
		require.NoError(t, err)
		require.Len(t, players, 1)
		require.Equal(t, players[0].Name, "Joe Bloggs")
	})

	t.Run("many expected", func(t *testing.T) {
		players, err := model.SearchPlayerByName(db, "JOE")
		require.NoError(t, err)
		require.Len(t, players, 3)
	})
}
