package search_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/db/search"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Dave")))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Chas")))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Danny")))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Greg")))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Sophie")))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Anne")))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Beryl")))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Charlotte")))

	t.Run("Default Search", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "dave")
		require.NoError(t, err)
		require.Len(t, players, 1)
		require.Equal(t, players[0].Name, "Dave")
	})

	t.Run("Default Search Fuzzy", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "dav")
		require.NoError(t, err)
		require.Len(t, players, 1)
		require.Equal(t, players[0].Name, "Dave")
	})
}
